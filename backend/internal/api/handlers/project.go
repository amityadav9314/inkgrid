package handlers

import (
	"github.com/amityadav9314/goinkgrid/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ProjectHandler handles project-related requests
type ProjectHandler struct {
	projectService services.ProjectService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// ProjectResponse represents a project response
type ProjectResponse struct {
	ID          uint           `json:"id"`
	UserID      uint           `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Settings    gin.H          `json:"settings"`
	Status      string         `json:"status"`
	MainImage   *ImageResponse `json:"main_image,omitempty"`
}

// CreateProjectRequest represents a create project request
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	MainImageID string `json:"main_image_id"`
	Settings    gin.H  `json:"settings"`
}

// UpdateProjectRequest represents an update project request
type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MainImageID string `json:"main_image_id"`
	Settings    gin.H  `json:"settings"`
}

// ListProjects returns a list of user projects
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// TODO: Get projects from database

	// Mock response for now
	c.JSON(http.StatusOK, gin.H{
		"projects": []ProjectResponse{
			{
				ID:          1,
				UserID:      userID.(uint),
				Name:        "Nature Mosaic",
				Description: "A mosaic made from nature photos",
				CreatedAt:   time.Now().Add(-24 * time.Hour),
				UpdatedAt:   time.Now().Add(-12 * time.Hour),
				Settings: gin.H{
					"tile_size":     50,
					"tile_density":  80,
					"overlay_ratio": 0.7,
					"style":         "classic",
				},
				Status: "completed",
				MainImage: &ImageResponse{
					ID:       "img-123",
					UserID:   userID.(uint),
					Type:     "main",
					Path:     "/uploads/img-123.jpg",
					Filename: "nature.jpg",
					Width:    1920,
					Height:   1080,
					Format:   "jpg",
				},
			},
			{
				ID:          2,
				UserID:      userID.(uint),
				Name:        "Family Collage",
				Description: "A mosaic of family photos",
				CreatedAt:   time.Now().Add(-48 * time.Hour),
				UpdatedAt:   time.Now().Add(-36 * time.Hour),
				Settings: gin.H{
					"tile_size":     30,
					"tile_density":  90,
					"overlay_ratio": 0.5,
					"style":         "flowing",
				},
				Status: "in_progress",
				MainImage: &ImageResponse{
					ID:       "img-456",
					UserID:   userID.(uint),
					Type:     "main",
					Path:     "/uploads/img-456.jpg",
					Filename: "family.jpg",
					Width:    2560,
					Height:   1440,
					Format:   "jpg",
				},
			},
		},
	})
}

// CreateProject creates a new project
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Validate that the main image belongs to the user
	// TODO: Create project in database

	// Mock response for now
	c.JSON(http.StatusCreated, ProjectResponse{
		ID:          3,
		UserID:      userID.(uint),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Settings:    req.Settings,
		Status:      "new",
	})
}

// GetProject returns a specific project
func (h *ProjectHandler) GetProject(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get project ID from path
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// TODO: Get project from database
	// TODO: Verify that the project belongs to the user

	// Mock response for now
	c.JSON(http.StatusOK, ProjectResponse{
		ID:          uint(projectID),
		UserID:      userID.(uint),
		Name:        "Nature Mosaic",
		Description: "A mosaic made from nature photos",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-12 * time.Hour),
		Settings: gin.H{
			"tile_size":     50,
			"tile_density":  80,
			"overlay_ratio": 0.7,
			"style":         "classic",
		},
		Status: "completed",
		MainImage: &ImageResponse{
			ID:       "img-123",
			UserID:   userID.(uint),
			Type:     "main",
			Path:     "/uploads/img-123.jpg",
			Filename: "nature.jpg",
			Width:    1920,
			Height:   1080,
			Format:   "jpg",
		},
	})
}

// UpdateProject updates a project
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get project ID from path
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get project from database
	// TODO: Verify that the project belongs to the user
	// TODO: Update project in database

	// Mock response for now
	c.JSON(http.StatusOK, ProjectResponse{
		ID:          uint(projectID),
		UserID:      userID.(uint),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
		Settings:    req.Settings,
		Status:      "updated",
	})
}

// DeleteProject deletes a project
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get project ID from path
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// TODO: Get project from database
	// TODO: Verify that the project belongs to the user
	// TODO: Delete project from database

	c.JSON(http.StatusOK, gin.H{
		"message": "Project deleted successfully",
		"id":      projectID,
		"user_id": userID,
	})
}
