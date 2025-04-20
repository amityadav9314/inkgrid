package db

import (
	"fmt"
	"github.com/amityadav9314/goinkgrid/internal/db/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// Try getting connection string from environment variable first
	// example:
	// export INKGRID_DATABASE_URL="host=localhost user=amityadav9314 password=amit8780 dbname=inkgrid port=5432 sslmode=disable"
	dsn := os.Getenv("INKGRID_DATABASE_URL")

	// If not set, use the provided credentials
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=localhost user=amityadav9314 password=amit8780 dbname=inkgrid port=5432 sslmode=disable",
		)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Auto migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Image{},
		&models.TileCollection{},
		&models.CollectionImage{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Successfully migrated database schema")
}
