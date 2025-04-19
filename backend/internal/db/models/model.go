package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Projects     []Project
}

type Project struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Settings    datatypes.JSON
	Status      string
	Images      []Image
}

type Image struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	ProjectID *uint  `gorm:"index"`
	Type      string // "main" or "tile"
	Path      string
	Filename  string
	Width     int
	Height    int
	Format    string
	CreatedAt time.Time
	ColorData datatypes.JSON // for tiles
}
