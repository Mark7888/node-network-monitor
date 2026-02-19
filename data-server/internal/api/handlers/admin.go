package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"mark7888/speedtest-data-server/internal/api/validators"
	"mark7888/speedtest-data-server/internal/auth"
	"mark7888/speedtest-data-server/internal/config"
	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AdminHandler handles admin-related endpoints
type AdminHandler struct {
	db         db.Database
	jwtManager *auth.JWTManager
	config     *config.Config
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(database db.Database, jwtManager *auth.JWTManager, cfg *config.Config) *AdminHandler {
	return &AdminHandler{
		db:         database,
		jwtManager: jwtManager,
		config:     cfg,
	}
}

// HandleLogin handles admin login
// POST /api/v1/admin/login
func (h *AdminHandler) HandleLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Verify credentials
	if req.Username != h.config.Admin.Username {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid credentials",
		})
		return
	}

	if !auth.VerifyPassword(req.Password, h.config.Admin.Password) {
		// Hash the password if it's not hashed yet (for first login)
		if req.Password == h.config.Admin.Password {
			// Direct match - password not hashed in config
			// In production, you should hash the password in config
		} else {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid credentials",
			})
			return
		}
	}

	// Generate JWT token
	token, expiresAt, err := h.jwtManager.Generate(req.Username)
	if err != nil {
		logger.Log.Error("Failed to generate JWT", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	logger.Log.Info("Admin login successful", zap.String("username", req.Username))

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Username:  req.Username,
	})
}

// HandleRefreshToken handles token refresh
// POST /api/v1/admin/refresh
func (h *AdminHandler) HandleRefreshToken(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	token, expiresAt, err := h.jwtManager.Generate(username.(string))
	if err != nil {
		logger.Log.Error("Failed to refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, models.RefreshTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	})
}

// HandleListNodes lists all nodes with optional filtering
// GET /api/v1/admin/nodes
func (h *AdminHandler) HandleListNodes(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	page, limit, _ = validators.ValidatePagination(page, limit)

	nodes, total, err := h.db.GetAllNodes(status, page, limit)
	if err != nil {
		logger.Log.Error("Failed to get nodes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve nodes",
		})
		return
	}

	// Enhance nodes with latest measurement info
	var nodesWithStats []models.NodeWithStats
	for _, node := range nodes {
		nodeStats, err := h.db.GetNodeWithStats(node.ID)
		if err != nil {
			// If we can't get stats, just use basic node info
			nodesWithStats = append(nodesWithStats, models.NodeWithStats{
				Node: node,
			})
			continue
		}
		nodesWithStats = append(nodesWithStats, *nodeStats)
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodesWithStats,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// HandleGetNodeDetails gets detailed information about a node
// GET /api/v1/admin/nodes/:id
func (h *AdminHandler) HandleGetNodeDetails(c *gin.Context) {
	nodeID, err := validators.ValidateUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid node ID",
		})
		return
	}

	nodeWithStats, err := h.db.GetNodeWithStats(nodeID)
	if err != nil {
		// Log at Info level for "not found" errors (expected 404s), Error level for actual DB issues
		if strings.Contains(err.Error(), "not found") {
			logger.Log.Info("Node not found", zap.String("node_id", nodeID.String()))
		} else {
			logger.Log.Error("Failed to get node details", zap.Error(err))
		}
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Node not found",
		})
		return
	}

	c.JSON(http.StatusOK, nodeWithStats)
}

// HandleGetNodeMeasurements gets measurements for a specific node
// GET /api/v1/admin/nodes/:id/measurements
func (h *AdminHandler) HandleGetNodeMeasurements(c *gin.Context) {
	nodeID, err := validators.ValidateUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid node ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1000"))
	page, limit, _ = validators.ValidatePagination(page, limit)

	// Parse status filter (all, successful, failed)
	status := c.DefaultQuery("status", "all")
	if status != "all" && status != "successful" && status != "failed" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid status parameter. Must be: all, successful, or failed",
		})
		return
	}

	// Parse time filters
	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	measurements, total, err := h.db.GetMeasurementsByNode(nodeID, from, to, page, limit, status)
	if err != nil {
		logger.Log.Error("Failed to get measurements", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve measurements",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"measurements": measurements,
		"total":        total,
		"page":         page,
		"limit":        limit,
	})
}

// HandleGetAggregatedMeasurements gets aggregated measurements for charting
// GET /api/v1/admin/measurements/aggregate
func (h *AdminHandler) HandleGetAggregatedMeasurements(c *gin.Context) {
	// Parse interval (required)
	interval := c.Query("interval")
	if interval == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "interval parameter is required",
		})
		return
	}

	// Validate interval
	if err := validators.ValidateInterval(interval); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Parse from time (required)
	fromStr := c.Query("from")
	if fromStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "from parameter is required",
		})
		return
	}
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid from timestamp",
			Details: err.Error(),
		})
		return
	}

	// Parse to time (required)
	toStr := c.Query("to")
	if toStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "to parameter is required",
		})
		return
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid to timestamp",
			Details: err.Error(),
		})
		return
	}

	// Validate time range
	if err := validators.ValidateTimeRange(&from, &to); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Parse node_ids (optional, can be multiple)
	var nodeIDs []uuid.UUID
	if nodeIDStrs := c.QueryArray("node_ids"); len(nodeIDStrs) > 0 {
		for _, idStr := range nodeIDStrs {
			id, err := validators.ValidateUUID(idStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "Invalid node ID",
					Details: err.Error(),
				})
				return
			}
			nodeIDs = append(nodeIDs, id)
		}
	}

	// Parse hide_archived (optional, defaults to false)
	hideArchived := false
	if hideArchivedStr := c.Query("hide_archived"); hideArchivedStr != "" {
		if hideArchivedStr == "true" {
			hideArchived = true
		}
	}

	measurements, err := h.db.GetAggregatedMeasurements(nodeIDs, from, to, interval, hideArchived)
	if err != nil {
		logger.Log.Error("Failed to get aggregated measurements", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve aggregated measurements",
		})
		return
	}

	c.JSON(http.StatusOK, models.AggregationResponse{
		Data:         measurements,
		Interval:     interval,
		TotalSamples: len(measurements),
	})
}

// HandleGetDashboard gets dashboard summary data
// GET /api/v1/admin/dashboard
func (h *AdminHandler) HandleGetDashboard(c *gin.Context) {
	// Get node counts
	totalNodes, activeNodes, unreachableNodes, inactiveNodes, err := h.db.GetNodeCounts()
	if err != nil {
		logger.Log.Error("Failed to get node counts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve dashboard data",
		})
		return
	}

	// Get measurement counts
	totalMeasurements, last24h, lastTimestamp, err := h.db.GetMeasurementCounts()
	if err != nil {
		logger.Log.Error("Failed to get measurement counts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve dashboard data",
		})
		return
	}

	// Get last 24h stats
	stats, err := h.db.GetLast24hStats()
	if err != nil {
		logger.Log.Error("Failed to get last 24h stats", zap.Error(err))
		stats = nil // Continue without stats
	}

	summary := models.DashboardSummary{
		TotalNodes:          totalNodes,
		ActiveNodes:         activeNodes,
		UnreachableNodes:    unreachableNodes,
		InactiveNodes:       inactiveNodes,
		TotalMeasurements:   totalMeasurements,
		MeasurementsLast24h: last24h,
		LastMeasurement:     lastTimestamp,
		AverageStats24h:     stats,
	}

	c.JSON(http.StatusOK, summary)
}

// HandleArchiveNode archives or unarchives a node
// PATCH /api/v1/admin/nodes/:id/archive
func (h *AdminHandler) HandleArchiveNode(c *gin.Context) {
	nodeID, err := validators.ValidateUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid node ID",
		})
		return
	}

	var req struct {
		Archived *bool `json:"archived" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if req.Archived == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "archived field is required",
		})
		return
	}

	err = h.db.ArchiveNode(nodeID, *req.Archived)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Node not found",
			})
			return
		}
		logger.Log.Error("Failed to archive node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update node",
		})
		return
	}

	action := "archived"
	if !*req.Archived {
		action = "unarchived"
	}

	logger.Log.Info("Node "+action,
		zap.String("node_id", nodeID.String()),
	)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Node " + action + " successfully",
		"archived": *req.Archived,
	})
}

// HandleSetNodeFavorite sets or removes the favorite status of a node
// PATCH /api/v1/admin/nodes/:id/favorite
func (h *AdminHandler) HandleSetNodeFavorite(c *gin.Context) {
	nodeID, err := validators.ValidateUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid node ID",
		})
		return
	}

	var req struct {
		Favorite *bool `json:"favorite" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if req.Favorite == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "favorite field is required",
		})
		return
	}

	err = h.db.SetNodeFavorite(nodeID, *req.Favorite)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Node not found",
			})
			return
		}
		logger.Log.Error("Failed to set node favorite", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update node",
		})
		return
	}

	action := "added to"
	if !*req.Favorite {
		action = "removed from"
	}

	logger.Log.Info("Node "+action+" favorites",
		zap.String("node_id", nodeID.String()),
	)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Node " + action + " favorites successfully",
		"favorite": *req.Favorite,
	})
}

// HandleDeleteNode deletes a node and all its measurements
// DELETE /api/v1/admin/nodes/:id
func (h *AdminHandler) HandleDeleteNode(c *gin.Context) {
	nodeID, err := validators.ValidateUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid node ID",
		})
		return
	}

	err = h.db.DeleteNode(nodeID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Node not found",
			})
			return
		}
		logger.Log.Error("Failed to delete node", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete node",
		})
		return
	}

	logger.Log.Info("Node deleted",
		zap.String("node_id", nodeID.String()),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Node deleted successfully",
	})
}
