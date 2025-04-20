package handlers

import (
	"fmt"
	"github.com/amityadav9314/goinkgrid/internal/db/models"
	"github.com/amityadav9314/goinkgrid/internal/services"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"image"
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

	// Get project ID from query or form data
	var projectID *uint

	// First try from query parameter
	projectIDStr := c.Query("project_id")

	// If not in query, try from form data
	if projectIDStr == "" {
		projectIDStr = c.PostForm("project_id")
	}

	if projectIDStr != "" {
		pid, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err == nil {
			pidUint := uint(pid)
			projectID = &pidUint
			fmt.Printf("Project ID parsed: %d\n", *projectID)
		} else {
			fmt.Printf("Error parsing project ID: %v\n", err)
		}
	} else {
		fmt.Println("No project ID provided")
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

	// Create directory structure: uploads/userID/projectID/
	userDir := filepath.Join(h.uploadDir, fmt.Sprintf("user_%d", userID.(uint)))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user directory"})
		return
	}

	var projectDir string
	if projectID != nil {
		projectDir = filepath.Join(userDir, fmt.Sprintf("project_%d", *projectID))
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project directory"})
			return
		}
	} else {
		projectDir = userDir
	}

	// Generate unique filename for main image
	filename := "main_" + uuid.New().String() + ext
	dst := filepath.Join(projectDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Get image dimensions using image package
	imgFile, err := os.Open(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}
	defer imgFile.Close()

	imgConfig, _, err := image.DecodeConfig(imgFile)
	var width, height int
	if err != nil {
		// If we can't decode the image, use placeholders
		width, height = 800, 600
	} else {
		width, height = imgConfig.Width, imgConfig.Height
	}

	// Create image metadata in database
	imageID := uuid.New().String()
	imagePath := strings.TrimPrefix(dst, h.uploadDir) // Store relative path
	if !strings.HasPrefix(imagePath, "/") {
		imagePath = "/" + imagePath
	}

	imageModel := &models.Image{
		UserID:    userID.(uint),
		ProjectID: projectID,
		Type:      "main",
		Path:      imagePath,
		Filename:  file.Filename,
		Width:     width,
		Height:    height,
		Format:    ext[1:],
	}

	// Save image metadata to database
	if err := h.imageService.Create(imageModel); err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Error saving image metadata: %v\n", err)
	} else {
		fmt.Printf("Image saved successfully with project ID: %v\n", projectID)
	}

	c.JSON(http.StatusOK, ImageResponse{
		ID:        imageID,
		UserID:    userID.(uint),
		ProjectID: projectID,
		Type:      "main",
		Path:      imagePath,
		Filename:  file.Filename,
		Width:     width,
		Height:    height,
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

	// Get project ID from query or form data
	var projectID *uint

	// First try from query parameter
	projectIDStr := c.Query("project_id")

	// If not in query, try from form data
	if projectIDStr == "" {
		projectIDStr = c.PostForm("project_id")
	}

	if projectIDStr != "" {
		pid, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err == nil {
			pidUint := uint(pid)
			projectID = &pidUint
			fmt.Printf("Tile upload - Project ID parsed: %d\n", *projectID)
		} else {
			fmt.Printf("Tile upload - Error parsing project ID: %v\n", err)
		}
	} else {
		fmt.Println("Tile upload - No project ID provided")
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

	// Create directory structure: uploads/userID/projectID/tiles/
	userDir := filepath.Join(h.uploadDir, fmt.Sprintf("user_%d", userID.(uint)))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user directory"})
		return
	}

	var projectDir string
	if projectID != nil {
		projectDir = filepath.Join(userDir, fmt.Sprintf("project_%d", *projectID))
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project directory"})
			return
		}
	} else {
		projectDir = userDir
	}

	// Create tiles directory
	tilesDir := filepath.Join(projectDir, "tiles")
	if err := os.MkdirAll(tilesDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tiles directory"})
		return
	}

	responses := make([]ImageResponse, 0, len(files))

	for i, file := range files {
		// Validate file type
		ext := filepath.Ext(file.Filename)
		if !isValidImageExt(ext) {
			continue // Skip invalid files
		}

		// Generate unique filename for tile image
		filename := fmt.Sprintf("tile_%d_%s%s", i+1, uuid.New().String(), ext)
		dst := filepath.Join(tilesDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, dst); err != nil {
			continue // Skip failed files
		}

		// Get image dimensions using image package
		imgFile, err := os.Open(dst)
		if err != nil {
			continue // Skip failed files
		}

		imgConfig, _, err := image.DecodeConfig(imgFile)
		imgFile.Close()

		var width, height int
		if err != nil {
			// If we can't decode the image, use placeholders
			width, height = 400, 300
		} else {
			width, height = imgConfig.Width, imgConfig.Height
		}

		// Create image metadata in database
		imageID := uuid.New().String()
		imagePath := strings.TrimPrefix(dst, h.uploadDir) // Store relative path
		if !strings.HasPrefix(imagePath, "/") {
			imagePath = "/" + imagePath
		}

		imageModel := &models.Image{
			UserID:    userID.(uint),
			ProjectID: projectID,
			Type:      "tile",
			Path:      imagePath,
			Filename:  file.Filename,
			Width:     width,
			Height:    height,
			Format:    ext[1:],
		}

		// Save image metadata to database
		if err := h.imageService.Create(imageModel); err != nil {
			// Log the error but don't fail the request
			fmt.Printf("Error saving tile image metadata: %v\n", err)
			continue
		} else {
			fmt.Printf("Tile image saved successfully with project ID: %v\n", projectID)
		}

		responses = append(responses, ImageResponse{
			ID:        imageID,
			UserID:    userID.(uint),
			ProjectID: projectID,
			Type:      "tile",
			Path:      imagePath,
			Filename:  file.Filename,
			Width:     width,
			Height:    height,
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

// GetProjectImages returns all images for a specific project
func (h *ImageHandler) GetProjectImages(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	//userID, exists := c.Get("userID")
	//if !exists {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	//	return
	//}

	// Get project ID from path parameter
	projectIDStr := c.Param("id")
	if projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Convert to uint
	projectIDUint := uint(projectID)

	// Get images from database
	images, err := h.imageService.FindByProjectID(projectIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	// Map database models to response format
	responses := make([]ImageResponse, 0, len(images))
	for _, img := range images {
		// Build full URL for image path
		imagePath := img.Path
		if !strings.HasPrefix(imagePath, "http") {
			// If it's a relative path, prepend the base URL
			baseURL := c.Request.Host
			scheme := "http"
			if c.Request.TLS != nil {
				scheme = "https"
			}
			imagePath = fmt.Sprintf("%s://%s%s", scheme, baseURL, imagePath)
		}

		responses = append(responses, ImageResponse{
			ID:        fmt.Sprintf("%d", img.ID), // Convert uint to string
			UserID:    img.UserID,
			ProjectID: img.ProjectID,
			Type:      img.Type,
			Path:      imagePath,
			Filename:  img.Filename,
			Width:     img.Width,
			Height:    img.Height,
			Format:    img.Format,
		})
	}

	// Group images by type
	mainImages := make([]ImageResponse, 0)
	tileImages := make([]ImageResponse, 0)

	for _, img := range responses {
		if img.Type == "main" {
			mainImages = append(mainImages, img)
		} else if img.Type == "tile" {
			tileImages = append(tileImages, img)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"main_images": mainImages,
		"tile_images": tileImages,
		"count":       len(responses),
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
