package handlers

import (
	"encoding/json"
	"github.com/amityadav9314/goinkgrid/internal/db/models"
	"github.com/amityadav9314/goinkgrid/internal/services"
	"gorm.io/datatypes"
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

	// Get projects from database
	projects, err := h.projectService.FindByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	// Convert to response format
	var projectResponses []ProjectResponse
	for _, project := range projects {
		// Create project response
		projectResponse := ProjectResponse{
			ID:          project.ID,
			UserID:      project.UserID,
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
			Settings:    gin.H{},
			Status:      project.Status,
		}

		// Parse settings if available
		if project.Settings != nil {
			var settings map[string]interface{}
			if err := json.Unmarshal(project.Settings, &settings); err == nil {
				projectResponse.Settings = settings
			}
		}

		// Add main image if available
		// TODO: In a real implementation, we would fetch the main image from the database
		// For now, we'll leave this empty

		projectResponses = append(projectResponses, projectResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projectResponses,
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

	// Create project in database
	project := models.Project{
		UserID:      userID.(uint),
		Name:        req.Name,
		Description: req.Description,
		Status:      "new",
	}

	// Convert settings to JSON if provided
	if req.Settings != nil {
		settingsJSON, err := json.Marshal(req.Settings)
		if err == nil {
			project.Settings = datatypes.JSON(settingsJSON)
		}
	}

	// Save project to database
	if err := h.projectService.Create(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Create response
	response := ProjectResponse{
		ID:          project.ID,
		UserID:      project.UserID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		Settings:    req.Settings,
		Status:      project.Status,
	}

	c.JSON(http.StatusCreated, response)
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

	// Get project from database
	project, err := h.projectService.FindByID(uint(projectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Verify that the project belongs to the user
	if project.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this project"})
		return
	}

	// Create response
	response := ProjectResponse{
		ID:          project.ID,
		UserID:      project.UserID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		Status:      project.Status,
	}

	// Parse settings if available
	if project.Settings != nil {
		var settings map[string]interface{}
		if err := json.Unmarshal(project.Settings, &settings); err == nil {
			response.Settings = settings
		}
	}

	// TODO: Fetch main image if needed

	c.JSON(http.StatusOK, response)
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

	// Get project from database
	project, err := h.projectService.FindByID(uint(projectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Verify that the project belongs to the user
	if project.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this project"})
		return
	}

	// Update project fields
	if req.Name != "" {
		project.Name = req.Name
	}

	// Description can be empty, so we update it regardless
	project.Description = req.Description

	// Update settings if provided
	if req.Settings != nil {
		settingsJSON, err := json.Marshal(req.Settings)
		if err == nil {
			project.Settings = datatypes.JSON(settingsJSON)
		}
	}

	// Update project in database
	if err := h.projectService.Update(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// Create response
	response := ProjectResponse{
		ID:          project.ID,
		UserID:      project.UserID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		Status:      project.Status,
	}

	// Parse settings
	if project.Settings != nil {
		var settings map[string]interface{}
		if err := json.Unmarshal(project.Settings, &settings); err == nil {
			response.Settings = settings
		}
	}

	c.JSON(http.StatusOK, response)
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

	// Delete project from database
	if err := h.projectService.Delete(uint(projectID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
