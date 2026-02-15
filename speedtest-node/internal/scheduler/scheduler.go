package scheduler

import (
	"mark7888/speedtest-node/internal/db"
	"mark7888/speedtest-node/internal/speedtest"
	"mark7888/speedtest-node/internal/sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler manages all scheduled tasks
type Scheduler struct {
	cron            *cron.Cron
	logger          *zap.Logger
	executor        *speedtest.Executor
	database        *db.DB
	sender          *sync.Sender
	aliveSender     *sync.AliveSender
	syncInterval    time.Duration
	aliveInterval   time.Duration
	retentionDays   int
	batchSize       int
	stopSyncChan    chan struct{}
	stopAliveChan   chan struct{}
	stopCleanupChan chan struct{}
}

// New creates a new scheduler
func New(
	speedtestCron string,
	executor *speedtest.Executor,
	database *db.DB,
	sender *sync.Sender,
	aliveSender *sync.AliveSender,
	syncInterval time.Duration,
	aliveInterval time.Duration,
	retentionDays int,
	batchSize int,
	logger *zap.Logger,
) (*Scheduler, error) {
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(&cronLogger{logger})))

	s := &Scheduler{
		cron:            c,
		logger:          logger,
		executor:        executor,
		database:        database,
		sender:          sender,
		aliveSender:     aliveSender,
		syncInterval:    syncInterval,
		aliveInterval:   aliveInterval,
		retentionDays:   retentionDays,
		batchSize:       batchSize,
		stopSyncChan:    make(chan struct{}),
		stopAliveChan:   make(chan struct{}),
		stopCleanupChan: make(chan struct{}),
	}

	// Schedule speedtest
	_, err := c.AddFunc(speedtestCron, s.runSpeedtest)
	if err != nil {
		return nil, err
	}

	logger.Info("Scheduler initialized", zap.String("speedtest_cron", speedtestCron))

	return s, nil
}

// Start starts all scheduled tasks
func (s *Scheduler) Start() {
	s.logger.Info("Starting scheduler")

	// Start cron scheduler (for speedtest)
	s.cron.Start()

	// Start background workers
	go s.syncWorker()
	go s.aliveWorker()
	go s.cleanupWorker()
}

// Stop stops all scheduled tasks
func (s *Scheduler) Stop() {
	s.logger.Info("Stopping scheduler")

	// Stop cron
	ctx := s.cron.Stop()
	<-ctx.Done()

	// Stop background workers
	close(s.stopSyncChan)
	close(s.stopAliveChan)
	close(s.stopCleanupChan)
}

// runSpeedtest executes a speedtest and stores the result
func (s *Scheduler) runSpeedtest() {
	s.logger.Info("Running scheduled speedtest")

	measurement, err := s.executor.Run()
	if err != nil {
		// Store as failed measurement
		s.database.InsertFailedMeasurement(time.Now().UTC(), err.Error(), 1)
		return
	}

	// Store measurement
	if err := s.database.InsertMeasurement(measurement); err != nil {
		s.logger.Error("Failed to store measurement", zap.Error(err))
	}
}

// syncWorker periodically syncs unsent measurements with the server
func (s *Scheduler) syncWorker() {
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.syncMeasurements()
			s.syncFailedMeasurements()
		case <-s.stopSyncChan:
			return
		}
	}
}

// syncMeasurements sends unsent measurements to the server
func (s *Scheduler) syncMeasurements() {
	// Skip if sender is not configured (offline mode)
	if s.sender == nil {
		return
	}

	measurements, err := s.database.GetUnsentMeasurements(s.batchSize)
	if err != nil {
		s.logger.Error("Failed to get unsent measurements", zap.Error(err))
		return
	}

	if len(measurements) == 0 {
		return
	}

	s.logger.Debug("Found unsent measurements", zap.Int("count", len(measurements)))

	if err := s.sender.SendMeasurements(measurements); err != nil {
		s.logger.Warn("Failed to sync measurements", zap.Error(err))
		return
	}

	// Mark as sent
	ids := make([]int64, len(measurements))
	for i, m := range measurements {
		ids[i] = m.ID
	}

	if err := s.database.MarkMeasurementsAsSent(ids); err != nil {
		s.logger.Error("Failed to mark measurements as sent", zap.Error(err))
	}
}

// syncFailedMeasurements sends unsent failed measurements to the server
func (s *Scheduler) syncFailedMeasurements() {
	// Skip if sender is not configured (offline mode)
	if s.sender == nil {
		return
	}

	failed, err := s.database.GetUnsentFailedMeasurements(s.batchSize)
	if err != nil {
		s.logger.Error("Failed to get unsent failed measurements", zap.Error(err))
		return
	}

	if len(failed) == 0 {
		return
	}

	s.logger.Debug("Found unsent failed measurements", zap.Int("count", len(failed)))

	if err := s.sender.SendFailedMeasurements(failed); err != nil {
		s.logger.Warn("Failed to sync failed measurements", zap.Error(err))
		return
	}

	// Mark as sent
	ids := make([]int64, len(failed))
	for i, f := range failed {
		ids[i] = f.ID
	}

	if err := s.database.MarkFailedMeasurementsAsSent(ids); err != nil {
		s.logger.Error("Failed to mark failed measurements as sent", zap.Error(err))
	}
}

// aliveWorker periodically sends alive signals to the server
func (s *Scheduler) aliveWorker() {
	// Skip if aliveSender is not configured (offline mode)
	if s.aliveSender == nil {
		<-s.stopAliveChan
		return
	}

	ticker := time.NewTicker(s.aliveInterval)
	defer ticker.Stop()

	// Send immediately on start
	if err := s.aliveSender.SendAlive(); err != nil {
		s.logger.Warn("Failed to send alive signal", zap.Error(err))
	}

	for {
		select {
		case <-ticker.C:
			if err := s.aliveSender.SendAlive(); err != nil {
				s.logger.Warn("Failed to send alive signal", zap.Error(err))
			}
		case <-s.stopAliveChan:
			return
		}
	}
}

// cleanupWorker periodically cleans up old data
func (s *Scheduler) cleanupWorker() {
	// Run cleanup once a day
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Run immediately on start
	s.cleanup()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.stopCleanupChan:
			return
		}
	}
}

// cleanup removes old data based on retention policy
func (s *Scheduler) cleanup() {
	retentionDate := time.Now().AddDate(0, 0, -s.retentionDays)

	s.logger.Info("Running cleanup", zap.Time("before", retentionDate))

	if err := s.database.DeleteMeasurementsBefore(retentionDate); err != nil {
		s.logger.Error("Failed to delete old measurements", zap.Error(err))
	}

	if err := s.database.DeleteFailedMeasurementsBefore(retentionDate); err != nil {
		s.logger.Error("Failed to delete old failed measurements", zap.Error(err))
	}
}

// cronLogger wraps zap logger for cron
type cronLogger struct {
	logger *zap.Logger
}

func (l *cronLogger) Printf(format string, v ...interface{}) {
	// Cron logs are too verbose, we'll skip them
	// If needed, uncomment the line below
	// l.logger.Debug(fmt.Sprintf(format, v...))
}
