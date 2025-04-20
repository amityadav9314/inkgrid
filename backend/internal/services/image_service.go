package services

import (
	"errors"
	"github.com/amityadav9314/goinkgrid/internal/db/models"
	"gorm.io/gorm"
)

// ImageServiceImpl implements the ImageService interface
type ImageServiceImpl struct {
	db *gorm.DB
}

// NewImageService creates a new ImageService implementation
func NewImageService(db *gorm.DB) ImageService {
	return &ImageServiceImpl{
		db: db,
	}
}

// FindByID finds an image by ID
func (s *ImageServiceImpl) FindByID(id uint) (*models.Image, error) {
	var image models.Image
	result := s.db.First(&image, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("image not found")
		}
		return nil, result.Error
	}
	return &image, nil
}

// FindByUserID finds images by user ID
func (s *ImageServiceImpl) FindByUserID(userID uint) ([]models.Image, error) {
	var images []models.Image
	result := s.db.Where("user_id = ?", userID).Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}
	return images, nil
}

// FindByProjectID finds images by project ID
func (s *ImageServiceImpl) FindByProjectID(projectID uint) ([]models.Image, error) {
	var images []models.Image
	result := s.db.Where("project_id = ?", projectID).Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}
	return images, nil
}

// Create creates a new image
func (s *ImageServiceImpl) Create(image *models.Image) error {
	result := s.db.Create(image)
	return result.Error
}

// Update updates an image
func (s *ImageServiceImpl) Update(image *models.Image) error {
	result := s.db.Save(image)
	return result.Error
}

// Delete deletes an image
func (s *ImageServiceImpl) Delete(id uint, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Image{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("image not found or you don't have permission to delete it")
	}
	return nil
}
