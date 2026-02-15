package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Initialize sets up the logger based on configuration
func Initialize(level, format, output string, outputConsole bool) error {
	// Parse log level
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	// Choose encoder format
	var encoder zapcore.Encoder
	if format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create output writers
	var writeSyncer zapcore.WriteSyncer
	if output == "" || output == "stdout" {
		writeSyncer = zapcore.AddSync(os.Stdout)
	} else {
		// Ensure log directory exists
		logDir := filepath.Dir(output)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		// Write to both file and stdout
		var writers []zapcore.WriteSyncer
		writers = append(writers, zapcore.AddSync(file))
		if outputConsole {
			writers = append(writers, zapcore.AddSync(os.Stdout))
		}
		writeSyncer = zapcore.NewMultiWriteSyncer(writers...)
	}

	// Create core and logger
	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
