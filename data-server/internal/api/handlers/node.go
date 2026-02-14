package handlers

import (
	"net/http"
	"time"

	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NodeHandler handles node-related endpoints
type NodeHandler struct {
	db *db.DB
}

// NewNodeHandler creates a new node handler
func NewNodeHandler(database *db.DB) *NodeHandler {
	return &NodeHandler{
		db: database,
	}
}

// HandleAlive handles node alive/registration signals
// POST /api/v1/node/alive
func (h *NodeHandler) HandleAlive(c *gin.Context) {
	var req models.AliveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid alive request", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Upsert node (create if doesn't exist, update if it does)
	err := h.db.UpsertNode(req.NodeID, req.NodeName)
	if err != nil {
		logger.Log.Error("Failed to upsert node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to register node",
		})
		return
	}

	logger.Log.Info("Node alive signal received",
		zap.String("node_id", req.NodeID.String()),
		zap.String("node_name", req.NodeName),
	)

	c.JSON(http.StatusOK, models.AliveResponse{
		Status:         "ok",
		ServerTime:     time.Now(),
		NodeRegistered: true,
	})
}
