package services

import (
	"errors"
	"github.com/amityadav9314/goinkgrid/internal/db/models"
	"gorm.io/gorm"
)

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	db *gorm.DB
}

// NewUserService creates a new UserService implementation
func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{
		db: db,
	}
}

// FindByEmail finds a user by email
func (s *UserServiceImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create creates a new user
func (s *UserServiceImpl) Create(user *models.User) error {
	result := s.db.Create(user)
	return result.Error
}
