package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MosaicHandler handles mosaic generation requests
type MosaicHandler struct {
	// TODO: Add mosaic service dependency
}

// NewMosaicHandler creates a new mosaic handler
func NewMosaicHandler() *MosaicHandler {
	return &MosaicHandler{}
}

// MosaicGenerationRequest represents the mosaic generation request
type MosaicGenerationRequest struct {
	ProjectID     *uint   `json:"project_id"`
	MainImageID   string  `json:"main_image_id" binding:"required"`
	TileImageIDs  []string `json:"tile_image_ids" binding:"required"`
	TileSize      int     `json:"tile_size" binding:"required,min=10,max=200"`
	TileDensity   int     `json:"tile_density" binding:"required,min=1,max=100"`
	OverlayRatio  float64 `json:"overlay_ratio" binding:"required,min=0,max=1"`
	Style         string  `json:"style" binding:"required,oneof=classic random flowing"`
	ColorCorrection bool   `json:"color_correction"`
}

// MosaicGenerationResponse represents the mosaic generation response
type MosaicGenerationResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// MosaicSettings represents the settings for mosaic generation
type MosaicSettings struct {
	TileSize        int     `json:"tile_size" binding:"required,min=10,max=200"`
	TileDensity     int     `json:"tile_density" binding:"required,min=1,max=100"`
	ColorAdjustment int     `json:"color_adjustment" binding:"required,min=0,max=100"`
	Style           string  `json:"style" binding:"required,oneof=classic random flowing"`
}

// GenerateMosaic handles mosaic generation requests
func (h *MosaicHandler) GenerateMosaic(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req MosaicGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Validate that the main image and tile images belong to the user
	// TODO: Create a new project if projectID is not provided
	// TODO: Start mosaic generation in a background goroutine

	// Generate a unique ID for the generation task
	generationID := uuid.New().String()

	// Return immediately with the generation ID
	c.JSON(http.StatusAccepted, MosaicGenerationResponse{
		ID:        generationID,
		Status:    "processing",
		CreatedAt: time.Now(),
	})

	// In a real implementation, we would start the generation in a background task
	// and update the status in the database as it progresses
}

// GetGenerationStatus returns the status of a mosaic generation task
func (h *MosaicHandler) GetGenerationStatus(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get generation ID from path
	generationID := c.Param("id")
	if generationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Generation ID is required"})
		return
	}

	// TODO: Get generation status from database
	// Mock response for now
	c.JSON(http.StatusOK, gin.H{
		"id":         generationID,
		"status":     "completed", // or "processing", "failed"
		"progress":   100,
		"created_at": time.Now().Add(-5 * time.Minute),
		"updated_at": time.Now(),
		"result_url": "/api/projects/1/mosaic.jpg", // If completed
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

	var settings MosaicSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Save settings to database associated with the user
	// For now, we'll just return success

	c.JSON(http.StatusOK, gin.H{
		"message": "Settings saved successfully",
		"user_id": userID,
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

	// TODO: Get settings from database for the user
	// For now, return default settings

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"settings": MosaicSettings{
			TileSize:        50,
			TileDensity:     80,
			ColorAdjustment: 50,
			Style:           "classic",
		},
	})
}
