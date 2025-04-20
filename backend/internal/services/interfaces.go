package services

import (
	"github.com/amityadav9314/goinkgrid/internal/db/models"
)

// UserService defines user-related operations
type UserService interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

// ProjectService defines project-related operations
type ProjectService interface {
	FindByID(id uint) (*models.Project, error)
	FindByUserID(userID uint) ([]models.Project, error)
	Create(project *models.Project) error
	Update(project *models.Project) error
	Delete(id uint, userID uint) error
}

// ImageService defines image-related operations
type ImageService interface {
	FindByID(id uint) (*models.Image, error)
	FindByUserID(userID uint) ([]models.Image, error)
	FindByProjectID(projectID uint) ([]models.Image, error)
	Create(image *models.Image) error
	Update(image *models.Image) error
	Delete(id uint, userID uint) error
	// Add other image-related methods
}

// MosaicService defines mosaic-related operations
type MosaicService interface {
	SaveSettings(userID uint, settings *models.MosaicSettings) error
	GetSettings(userID uint, projectID *uint) (*models.MosaicSettings, error)
	GenerateMosaic(userID uint, projectID uint, mainImageID uint, tileImageIDs []uint, settings *models.MosaicSettings) (*models.GeneratedMosaic, error)
	GetMosaicStatus(userID uint, mosaicID uint) (*models.GeneratedMosaic, error)
	GetProjectMosaics(userID uint, projectID uint) ([]models.GeneratedMosaic, error)
}
