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
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index"`
	ProjectID   *uint  `gorm:"index"`
	Type        string // "main" or "tile"
	Path        string
	Filename    string
	Width       int
	Height      int
	Format      string
	CreatedAt   time.Time
	ColorData   datatypes.JSON   // for tiles
	Collections []TileCollection `gorm:"many2many:collection_images;"`
}

// TileCollection represents a group of tile images
type TileCollection struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	ProjectID *uint  `gorm:"index"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	Images    []Image `gorm:"many2many:collection_images;"`
}

// CollectionImage represents the many-to-many relationship
// between collections and images
type CollectionImage struct {
	CollectionID uint `gorm:"primaryKey"`
	ImageID      uint `gorm:"primaryKey"`
}

// MosaicSettings represents user-specific mosaic generation settings
type MosaicSettings struct {
	ID              uint   `gorm:"primaryKey"`
	UserID          uint   `gorm:"not null;index;uniqueIndex"` // One settings per user
	ProjectID       *uint  `gorm:"index"`
	TileSize        int    `gorm:"not null;default:50"`
	TileDensity     int    `gorm:"not null;default:80"`
	ColorAdjustment int    `gorm:"not null;default:50"`
	Style           string `gorm:"not null;default:'classic'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Update the existing Image model to add the collections relationship
func init() {
}
