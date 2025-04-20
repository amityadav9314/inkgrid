package services

import (
	"github.com/amityadav9314/goinkgrid/internal/db/models"
	"github.com/amityadav9314/goinkgrid/internal/db/postgres"
)

// MosaicService handles mosaic-related operations
type MosaicService interface {
	SaveSettings(userID uint, settings *models.MosaicSettings) error
	GetSettings(userID uint) (*models.MosaicSettings, error)
}

type mosaicService struct{}

// NewMosaicService creates a new mosaic service
func NewMosaicService() MosaicService {
	return &mosaicService{}
}

// SaveSettings saves or updates mosaic settings for a user
func (s *mosaicService) SaveSettings(userID uint, settings *models.MosaicSettings) error {
	// Check if settings already exist for this user
	var existingSettings models.MosaicSettings
	result := db.DB.Where("user_id = ? and project_id = ?", userID, settings.ProjectID).First(&existingSettings)

	if result.Error == nil {
		// Update existing settings
		existingSettings.TileSize = settings.TileSize
		existingSettings.TileDensity = settings.TileDensity
		existingSettings.ColorAdjustment = settings.ColorAdjustment
		existingSettings.Style = settings.Style
		return db.DB.Save(&existingSettings).Error
	}

	// Create new settings
	settings.UserID = userID
	return db.DB.Create(settings).Error
}

// GetSettings retrieves mosaic settings for a user
func (s *mosaicService) GetSettings(userID uint) (*models.MosaicSettings, error) {
	var settings models.MosaicSettings
	result := db.DB.Where("user_id = ?", userID).First(&settings)

	if result.Error != nil {
		// If no settings found, return default settings
		return &models.MosaicSettings{
			UserID:          userID,
			TileSize:        50,
			TileDensity:     80,
			ColorAdjustment: 50,
			Style:           "classic",
		}, nil
	}

	return &settings, nil
}
