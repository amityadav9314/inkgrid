package services

import (
	"errors"
	"github.com/amityadav9314/goinkgrid/internal/db/models"
	"gorm.io/gorm"
)

type ProjectServiceImpl struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) ProjectService {
	return &ProjectServiceImpl{
		db: db,
	}
}

func (s *ProjectServiceImpl) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	result := s.db.First(&project, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, result.Error
	}
	return &project, nil
}

func (s *ProjectServiceImpl) FindByUserID(userID uint) ([]models.Project, error) {
	var projects []models.Project
	result := s.db.Where("user_id = ?", userID).Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (s *ProjectServiceImpl) Create(project *models.Project) error {
	result := s.db.Create(project)
	return result.Error
}

func (s *ProjectServiceImpl) Update(project *models.Project) error {
	result := s.db.Save(project)
	return result.Error
}

func (s *ProjectServiceImpl) Delete(id uint, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Project{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("project not found or you don't have permission to delete it")
	}
	return nil
}
