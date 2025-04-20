package handlers

import (
	"github.com/amityadav9314/goinkgrid/internal/services"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ImageHandler handles image-related requests
type ImageHandler struct {
	uploadDir    string
	imageService services.ImageService
}

// NewImageHandler creates a new image handler
func NewImageHandler(uploadPath string, imageService services.ImageService) *ImageHandler {
	return &ImageHandler{
		uploadDir:    uploadPath,
		imageService: imageService,
	}
}

// ImageResponse represents an image response
type ImageResponse struct {
	ID        string `json:"id"`
	UserID    uint   `json:"user_id"`
	ProjectID *uint  `json:"project_id,omitempty"`
	Type      string `json:"type"` // main or tile
	Path      string `json:"path"`
	Filename  string `json:"filename"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
}

// UploadMainImage handles main image upload
func (h *ImageHandler) UploadMainImage(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get project ID from query if present
	var projectID *uint
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		pid := uint(0)
		// TODO: Convert projectIDStr to uint
		projectID = &pid
	}

	// Get file from form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	if !isValidImageExt(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + ext
	dst := filepath.Join(h.uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// TODO: Process image (get dimensions, format, etc.)
	// TODO: Save image metadata to database

	c.JSON(http.StatusOK, ImageResponse{
		ID:        uuid.New().String(),
		UserID:    userID.(uint),
		ProjectID: projectID,
		Type:      "main",
		Path:      dst,
		Filename:  file.Filename,
		Width:     800, // Placeholder
		Height:    600, // Placeholder
		Format:    ext[1:],
	})
}

// UploadTileImages handles batch upload of tile images
func (h *ImageHandler) UploadTileImages(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get project ID from query if present
	var projectID *uint
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		pid := uint(0)
		// TODO: Convert projectIDStr to uint
		projectID = &pid
	}

	// Get collection ID from query if present
	// Commented out until implemented
	/*
		var collectionID *uint
		if collectionIDStr := c.Query("collection_id"); collectionIDStr != "" {
			cid := uint(0)
			// TODO: Convert collectionIDStr to uint
			collectionID = &cid
		}
	*/

	// Get files from form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	files := form.File["images[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images provided"})
		return
	}

	responses := make([]ImageResponse, 0, len(files))

	for _, file := range files {
		// Validate file type
		ext := filepath.Ext(file.Filename)
		if !isValidImageExt(ext) {
			continue // Skip invalid files
		}

		// Generate unique filename
		filename := uuid.New().String() + ext
		dst := filepath.Join(h.uploadDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, dst); err != nil {
			continue // Skip failed files
		}

		// TODO: Process image (get dimensions, format, etc.)
		// TODO: Save image metadata to database
		// TODO: If collectionID is provided, add image to collection

		responses = append(responses, ImageResponse{
			ID:        uuid.New().String(),
			UserID:    userID.(uint),
			ProjectID: projectID,
			Type:      "tile",
			Path:      dst,
			Filename:  file.Filename,
			Width:     400, // Placeholder
			Height:    300, // Placeholder
			Format:    ext[1:],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Images uploaded successfully",
		"count":   len(responses),
		"images":  responses,
	})
}

// GetTileCollections returns the user's tile collections
func (h *ImageHandler) GetTileCollections(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// TODO: Get collections from database

	// Mock response for now
	c.JSON(http.StatusOK, gin.H{
		"collections": []gin.H{
			{
				"id":          1,
				"user_id":     userID,
				"name":        "Nature",
				"created_at":  "2023-01-01T00:00:00Z",
				"image_count": 10,
			},
			{
				"id":          2,
				"user_id":     userID,
				"name":        "Family",
				"created_at":  "2023-01-02T00:00:00Z",
				"image_count": 15,
			},
		},
	})
}

// Helper function to validate image extensions
func isValidImageExt(ext string) bool {
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".heic": true,
	}
	return validExts[ext]
}
