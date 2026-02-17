package handlers

import (
	"net/http"

	"mark7888/speedtest-data-server/internal/api/validators"
	"mark7888/speedtest-data-server/internal/auth"
	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// APIKeyHandler handles API key management endpoints
type APIKeyHandler struct {
	db db.Database
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(database db.Database) *APIKeyHandler {
	return &APIKeyHandler{
		db: database,
	}
}

// HandleListAPIKeys lists all API keys
// GET /api/v1/admin/api-keys
func (h *APIKeyHandler) HandleListAPIKeys(c *gin.Context) {
	apiKeys, err := h.db.GetAllAPIKeys()
	if err != nil {
		logger.Log.Error("Failed to get API keys", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve API keys",
		})
		return
	}

	c.JSON(http.StatusOK, models.ListAPIKeysResponse{
		APIKeys: apiKeys,
		Total:   len(apiKeys),
	})
}

// HandleCreateAPIKey creates a new API key
// POST /api/v1/admin/api-keys
func (h *APIKeyHandler) HandleCreateAPIKey(c *gin.Context) {
	var req models.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get username from context
	username, _ := c.Get("username")
	createdBy := ""
	if u, ok := username.(string); ok {
		createdBy = u
	}

	// Generate API key
	plainKey, err := auth.GenerateAPIKey()
	if err != nil {
		logger.Log.Error("Failed to generate API key", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate API key",
		})
		return
	}

	// Create API key in database
	apiKey, err := h.db.CreateAPIKey(req.Name, plainKey, createdBy)
	if err != nil {
		logger.Log.Error("Failed to create API key", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create API key",
		})
		return
	}

	logger.Log.Info("API key created",
		zap.String("id", apiKey.ID.String()),
		zap.String("name", req.Name),
		zap.String("created_by", createdBy),
	)

	c.JSON(http.StatusCreated, models.CreateAPIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Key:       plainKey, // Only shown once
		Enabled:   apiKey.Enabled,
		CreatedAt: apiKey.CreatedAt,
		Warning:   "Save this key securely. It won't be shown again.",
	})
}

// HandleUpdateAPIKey updates an API key (enable/disable)
// PATCH /api/v1/admin/api-keys/:id
func (h *APIKeyHandler) HandleUpdateAPIKey(c *gin.Context) {
	keyID, err := validators.ValidateUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid API key ID",
		})
		return
	}

	var req models.UpdateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if req.Enabled != nil {
		err := h.db.UpdateAPIKeyEnabled(keyID, *req.Enabled)
		if err != nil {
			logger.Log.Error("Failed to update API key", zap.Error(err))
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to update API key",
			})
			return
		}
	}

	// Get updated API key
	apiKey, err := h.db.GetAPIKeyByID(keyID)
	if err != nil {
		logger.Log.Error("Failed to get updated API key", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve updated API key",
		})
		return
	}

	logger.Log.Info("API key updated",
		zap.String("id", keyID.String()),
		zap.Bool("enabled", apiKey.Enabled),
	)

	c.JSON(http.StatusOK, apiKey)
}

// HandleDeleteAPIKey deletes an API key
// DELETE /api/v1/admin/api-keys/:id
func (h *APIKeyHandler) HandleDeleteAPIKey(c *gin.Context) {
	keyID, err := validators.ValidateUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid API key ID",
		})
		return
	}

	err = h.db.DeleteAPIKey(keyID)
	if err != nil {
		logger.Log.Error("Failed to delete API key", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete API key",
		})
		return
	}

	logger.Log.Info("API key deleted", zap.String("id", keyID.String()))

	c.Status(http.StatusNoContent)
}
