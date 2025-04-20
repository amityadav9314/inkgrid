package handlers

import (
	"net/http"
	"strconv"
	"time"

	"fmt"
	"github.com/amityadav9314/goinkgrid/internal/db/models"
	"github.com/amityadav9314/goinkgrid/internal/services"
	"github.com/gin-gonic/gin"
)

// MosaicHandler handles mosaic generation requests
type MosaicHandler struct {
	mosaicService services.MosaicService
}

// NewMosaicHandler creates a new mosaic handler
func NewMosaicHandler(mosaicService services.MosaicService) *MosaicHandler {
	return &MosaicHandler{
		mosaicService: mosaicService,
	}
}

// MosaicGenerationRequest represents the mosaic generation request
type MosaicGenerationRequest struct {
	ProjectID       *uint    `json:"project_id"`
	MainImageID     string   `json:"main_image_id" binding:"required"`
	TileImageIDs    []string `json:"tile_image_ids" binding:"required"`
	TileSize        int      `json:"tile_size" binding:"required,min=10,max=200"`
	TileDensity     int      `json:"tile_density" binding:"required,min=1,max=100"`
	OverlayRatio    float64  `json:"overlay_ratio" binding:"required,min=0,max=1"`
	Style           string   `json:"style" binding:"required,oneof=classic random flowing"`
	ColorCorrection bool     `json:"color_correction"`
}

// MosaicGenerationResponse represents the mosaic generation response
type MosaicGenerationResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// MosaicSettings represents the settings for mosaic generation
type MosaicSettings struct {
	TileSize        int    `json:"tile_size" binding:"required,min=10,max=200"`
	TileDensity     int    `json:"tile_density" binding:"required,min=1,max=100"`
	ColorAdjustment int    `json:"color_adjustment" binding:"required,min=0,max=100"`
	Style           string `json:"style" binding:"required,oneof=classic random flowing"`
}

// GenerateMosaic handles mosaic generation requests
func (h *MosaicHandler) GenerateMosaic(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req MosaicGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that project ID is provided
	if req.ProjectID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	// Parse main image ID
	mainImageID, err := strconv.ParseUint(req.MainImageID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid main image ID"})
		return
	}

	// Parse tile image IDs
	tileImageIDs := make([]uint, 0, len(req.TileImageIDs))
	for _, idStr := range req.TileImageIDs {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			continue
		}
		tileImageIDs = append(tileImageIDs, uint(id))
	}

	if len(tileImageIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid tile image IDs provided"})
		return
	}

	// Create settings for mosaic generation
	settings := &models.MosaicSettings{
		UserID:          userID.(uint),
		ProjectID:       req.ProjectID,
		TileSize:        req.TileSize,
		TileDensity:     req.TileDensity,
		ColorAdjustment: int(req.OverlayRatio * 100),
		Style:           req.Style,
	}

	// Start mosaic generation
	mosaic, err := h.mosaicService.GenerateMosaic(
		userID.(uint),
		*req.ProjectID,
		uint(mainImageID),
		tileImageIDs,
		settings,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the mosaic generation ID
	c.JSON(http.StatusAccepted, MosaicGenerationResponse{
		ID:        fmt.Sprintf("%d", mosaic.ID),
		Status:    mosaic.Status,
		CreatedAt: mosaic.CreatedAt,
	})
}

// GetGenerationStatus returns the status of a mosaic generation task
func (h *MosaicHandler) GetGenerationStatus(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get generation ID from path
	generationIDStr := c.Param("id")
	if generationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Generation ID is required"})
		return
	}

	// Parse generation ID
	generationID, err := strconv.ParseUint(generationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid generation ID"})
		return
	}

	// Get mosaic status
	mosaic, err := h.mosaicService.GetMosaicStatus(userID.(uint), uint(generationID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := gin.H{
		"id":         fmt.Sprintf("%d", mosaic.ID),
		"status":     mosaic.Status,
		"progress":   mosaic.Progress,
		"created_at": mosaic.CreatedAt,
		"updated_at": mosaic.UpdatedAt,
	}

	// Add result URLs if completed
	if mosaic.Status == "completed" {
		// Build full URLs for the mosaic images
		baseURL := c.Request.Host
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}

		sdURL := fmt.Sprintf("%s://%s/uploads%s", scheme, baseURL, mosaic.SDPath)
		hdURL := fmt.Sprintf("%s://%s/uploads%s", scheme, baseURL, mosaic.HDPath)

		response["sd_url"] = sdURL
		response["hd_url"] = hdURL
	}

	// Add error message if failed
	if mosaic.Status == "failed" && mosaic.ErrorMessage != "" {
		response["error"] = mosaic.ErrorMessage
	}

	c.JSON(http.StatusOK, response)
}

// GetProjectMosaics returns all mosaics for a project
func (h *MosaicHandler) GetProjectMosaics(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get project ID from path
	projectIDStr := c.Param("id")
	if projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	// Parse project ID
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Get mosaics for the project
	mosaics, err := h.mosaicService.GetProjectMosaics(userID.(uint), uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := make([]gin.H, 0, len(mosaics))
	for _, mosaic := range mosaics {
		// Build full URLs for the mosaic images
		baseURL := c.Request.Host
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}

		mosaicResponse := gin.H{
			"id":               fmt.Sprintf("%d", mosaic.ID),
			"status":           mosaic.Status,
			"progress":         mosaic.Progress,
			"created_at":       mosaic.CreatedAt,
			"updated_at":       mosaic.UpdatedAt,
			"tile_size":        mosaic.TileSize,
			"tile_density":     mosaic.TileDensity,
			"color_adjustment": mosaic.ColorAdjustment,
			"style":            mosaic.Style,
		}

		// Add URLs if completed
		if mosaic.Status == "completed" {
			if mosaic.SDPath != "" {
				mosaicResponse["sd_url"] = fmt.Sprintf("%s://%s%s", scheme, baseURL, mosaic.SDPath)
			}
			if mosaic.HDPath != "" {
				mosaicResponse["hd_url"] = fmt.Sprintf("%s://%s%s", scheme, baseURL, mosaic.HDPath)
			}
		}

		// Add error message if failed
		if mosaic.Status == "failed" && mosaic.ErrorMessage != "" {
			mosaicResponse["error"] = mosaic.ErrorMessage
		}

		response = append(response, mosaicResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"mosaics": response,
		"count":   len(response),
	})
}

// SaveMosaicSettings handles saving mosaic settings
func (h *MosaicHandler) SaveMosaicSettings(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var requestBody struct {
		TileSize        int    `json:"tile_size" binding:"required,min=10,max=200"`
		TileDensity     int    `json:"tile_density" binding:"required,min=1,max=100"`
		ColorAdjustment int    `json:"color_adjustment" binding:"required,min=0,max=100"`
		Style           string `json:"style" binding:"required,oneof=classic random flowing"`
		ProjectID       *uint  `json:"project_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create settings model
	settings := &models.MosaicSettings{
		UserID:          userID.(uint),
		ProjectID:       requestBody.ProjectID,
		TileSize:        requestBody.TileSize,
		TileDensity:     requestBody.TileDensity,
		ColorAdjustment: requestBody.ColorAdjustment,
		Style:           requestBody.Style,
	}

	// Save settings to database associated with the user
	err := h.mosaicService.SaveSettings(userID.(uint), settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Settings saved successfully",
		"user_id":  userID,
		"settings": settings,
	})
}

// GetMosaicSettings returns the saved mosaic settings for a user
func (h *MosaicHandler) GetMosaicSettings(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if project ID is provided in query parameters
	var projectID *uint
	projectIDStr := c.Query("project_id")
	if projectIDStr != "" {
		// Convert string to uint
		pid, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err == nil {
			pidUint := uint(pid)
			projectID = &pidUint
		}
	}

	// Get settings from database for the user and project
	settings, err := h.mosaicService.GetSettings(userID.(uint), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"settings": settings,
	})
}
