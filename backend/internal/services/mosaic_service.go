package services

import (
	"errors"
	"fmt"
	models "github.com/amityadav9314/goinkgrid/internal/db/models"
	db "github.com/amityadav9314/goinkgrid/internal/db/postgres"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	xdraw "golang.org/x/image/draw"
	"gorm.io/gorm"
)

type MosaicServiceImpl struct {
	uploadDir string
	// Map to track active generation tasks
	activeTasks     map[uint]bool
	activeTasksLock sync.Mutex
}

// NewMosaicService creates a new mosaic service
func NewMosaicService(uploadDir string) MosaicService {
	return &MosaicServiceImpl{
		uploadDir:   uploadDir,
		activeTasks: make(map[uint]bool),
	}
}

// SaveSettings saves or updates mosaic settings for a user
func (s *MosaicServiceImpl) SaveSettings(userID uint, settings *models.MosaicSettings) error {
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
func (s *MosaicServiceImpl) GetSettings(userID uint, projectID *uint) (*models.MosaicSettings, error) {
	var settings models.MosaicSettings

	// Build query based on whether projectID is provided
	query := db.DB.Where("user_id = ?", userID)
	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	}

	result := query.First(&settings)

	if result.Error != nil {
		// If no settings found, return default settings
		defaultSettings := &models.MosaicSettings{
			UserID:          userID,
			TileSize:        50,
			TileDensity:     80,
			ColorAdjustment: 50,
			Style:           "classic",
		}

		// Set project ID if provided
		if projectID != nil {
			defaultSettings.ProjectID = projectID
		}

		return defaultSettings, nil
	}

	return &settings, nil
}

// GenerateMosaic generates a mosaic image from a main image and tile images
func (s *MosaicServiceImpl) GenerateMosaic(userID uint, projectID uint, mainImageID uint, tileImageIDs []uint, settings *models.MosaicSettings) (*models.GeneratedMosaic, error) {
	// Check if we already have an active task for this project
	s.activeTasksLock.Lock()
	if s.activeTasks[projectID] {
		s.activeTasksLock.Unlock()
		return nil, errors.New("a mosaic generation is already in progress for this project")
	}
	s.activeTasks[projectID] = true
	s.activeTasksLock.Unlock()

	// Create a new GeneratedMosaic record
	mosaic := &models.GeneratedMosaic{
		UserID:          userID,
		ProjectID:       projectID,
		MainImageID:     mainImageID,
		Status:          "processing",
		TileSize:        settings.TileSize,
		TileDensity:     settings.TileDensity,
		ColorAdjustment: settings.ColorAdjustment,
		Style:           settings.Style,
		Progress:        0,
	}

	// Save the initial record
	if err := db.DB.Create(mosaic).Error; err != nil {
		s.activeTasksLock.Lock()
		delete(s.activeTasks, projectID)
		s.activeTasksLock.Unlock()
		return nil, err
	}

	// Start the generation process in a goroutine
	go s.generateMosaicAsync(mosaic, mainImageID, tileImageIDs)

	return mosaic, nil
}

// generateMosaicAsync handles the asynchronous mosaic generation process
func (s *MosaicServiceImpl) generateMosaicAsync(mosaic *models.GeneratedMosaic, mainImageID uint, tileImageIDs []uint) {
	defer func() {
		// Remove from active tasks when done
		s.activeTasksLock.Lock()
		delete(s.activeTasks, mosaic.ProjectID)
		s.activeTasksLock.Unlock()
	}()

	// Update progress to 10%
	mosaic.Progress = 10
	db.DB.Save(mosaic)

	// Get the main image
	var mainImage models.Image
	if err := db.DB.First(&mainImage, mainImageID).Error; err != nil {
		mosaic.Status = "failed"
		mosaic.ErrorMessage = "Failed to find main image"
		db.DB.Save(mosaic)
		return
	}

	// Get the tile images
	var tileImages []models.Image
	if err := db.DB.Where("id IN ?", tileImageIDs).Find(&tileImages).Error; err != nil {
		mosaic.Status = "failed"
		mosaic.ErrorMessage = "Failed to find tile images"
		db.DB.Save(mosaic)
		return
	}

	// Update progress to 20%
	mosaic.Progress = 20
	db.DB.Save(mosaic)

	// Create directory for the mosaic
	mosaicDir := filepath.Join(s.uploadDir, fmt.Sprintf("user_%d", mosaic.UserID), fmt.Sprintf("project_%d", mosaic.ProjectID), "mosaics")
	if err := os.MkdirAll(mosaicDir, 0755); err != nil {
		mosaic.Status = "failed"
		mosaic.ErrorMessage = "Failed to create mosaic directory"
		db.DB.Save(mosaic)
		return
	}

	// Generate unique filenames for SD and HD mosaics
	timestamp := time.Now().Format("20060102150405")
	sdFilename := fmt.Sprintf("mosaic_sd_%s.jpg", timestamp)
	hdFilename := fmt.Sprintf("mosaic_hd_%s.jpg", timestamp)
	sdPath := filepath.Join(mosaicDir, sdFilename)
	hdPath := filepath.Join(mosaicDir, hdFilename)

	// Update progress to 30%
	mosaic.Progress = 30
	db.DB.Save(mosaic)

	// Simulate mosaic generation (in a real implementation, this would be the actual generation code)
	// For now, we'll just create placeholder images
	if err := s.createPlaceholderMosaics(mainImage.Path, tileImages, sdPath, hdPath, mosaic); err != nil {
		mosaic.Status = "failed"
		mosaic.ErrorMessage = fmt.Sprintf("Failed to generate mosaic: %v", err)
		db.DB.Save(mosaic)
		return
	}

	// Save the paths to the generated mosaics
	mosaic.SDPath = strings.TrimPrefix(sdPath, s.uploadDir)
	if !strings.HasPrefix(mosaic.SDPath, "/") {
		mosaic.SDPath = "/" + mosaic.SDPath
	}

	mosaic.HDPath = strings.TrimPrefix(hdPath, s.uploadDir)
	if !strings.HasPrefix(mosaic.HDPath, "/") {
		mosaic.HDPath = "/" + mosaic.HDPath
	}

	mosaic.Status = "completed"
	mosaic.Progress = 100
	db.DB.Save(mosaic)
}

// createPlaceholderMosaics creates placeholder mosaic images for development
// In a real implementation, this would be replaced with actual mosaic generation logic
func (s *MosaicServiceImpl) createPlaceholderMosaics(mainImagePath string, tileImages []models.Image, sdPath, hdPath string, mosaic *models.GeneratedMosaic) error {
	// Open the main image
	// Get the absolute path to the project root directory
	projectRoot, _ := filepath.Abs(".")

	// Construct the path to the main image based on the actual directory structure
	// The uploadDir is "./uploads" but the actual path is "/backend/uploads"
	var fullMainImagePath string

	// Check if mainImagePath already has the correct structure
	if strings.HasPrefix(mainImagePath, "/uploads/") {
		// Path is already in the correct format, just need to resolve from project root
		fullMainImagePath = filepath.Join(projectRoot, mainImagePath)
	} else if strings.HasPrefix(mainImagePath, "uploads/") {
		// Path starts with "uploads/" without leading slash
		fullMainImagePath = filepath.Join(projectRoot, "/", mainImagePath)
	} else {
		// Assume it's a relative path to the uploads directory
		fullMainImagePath = filepath.Join(projectRoot, "backend/uploads", strings.TrimPrefix(mainImagePath, "/"))
	}

	// Add logging to debug path issues
	fmt.Printf("Opening main image at path: %s\n", fullMainImagePath)
	fmt.Printf("Project root: %s\n", projectRoot)
	fmt.Printf("Original main image path: %s\n", mainImagePath)

	// Check if file exists before attempting to open
	if _, err := os.Stat(fullMainImagePath); os.IsNotExist(err) {
		// Try an alternative approach - use absolute path to backend/uploads
		backendUploadsPath, _ := filepath.Abs(filepath.Join(projectRoot, "uploads"))
		alternativePath := filepath.Join(backendUploadsPath, strings.TrimPrefix(mainImagePath, "/uploads/"))
		alternativePath = strings.TrimPrefix(alternativePath, "/uploads/")

		fmt.Printf("Original path not found, trying alternative path: %s\n", alternativePath)

		if _, err := os.Stat(alternativePath); os.IsNotExist(err) {
			// One more attempt - try to find the file directly in the file system
			findCmd := fmt.Sprintf("find %s -name '%s' 2>/dev/null",
				filepath.Join(projectRoot, "backend"),
				filepath.Base(mainImagePath))

			fmt.Printf("Running find command: %s\n", findCmd)

			cmd := exec.Command("sh", "-c", findCmd)
			output, _ := cmd.Output()

			if len(output) > 0 {
				foundPath := strings.TrimSpace(string(output))
				fmt.Printf("Found file at: %s\n", foundPath)
				fullMainImagePath = foundPath
			} else {
				return fmt.Errorf("failed to find main image: %v (tried paths: %s and %s)",
					err, fullMainImagePath, alternativePath)
			}
		} else {
			fullMainImagePath = alternativePath
		}
	}

	mainImg, err := openImage(fullMainImagePath)
	if err != nil {
		return fmt.Errorf("failed to open main image: %v", err)
	}

	// Update progress to 40%
	mosaic.Progress = 40
	db.DB.Save(mosaic)

	// Get the dimensions of the main image
	bounds := mainImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create SD and HD images
	sdWidth, sdHeight := width/2, height/2
	hdWidth, hdHeight := width, height

	// Create new images
	sdImg := image.NewRGBA(image.Rect(0, 0, sdWidth, sdHeight))
	hdImg := image.NewRGBA(image.Rect(0, 0, hdWidth, hdHeight))

	// Update progress to 50%
	mosaic.Progress = 50
	db.DB.Save(mosaic)

	// Draw a grid of tiles on the images
	tileSize := mosaic.TileSize
	rand.Seed(time.Now().UnixNano())

	// Load tile images
	tileImgList := make([]image.Image, 0, len(tileImages))
	for _, tile := range tileImages {
		fullTilePath := filepath.Join(projectRoot, strings.TrimPrefix(tile.Path, "/"))
		tileImg, err := openImage(fullTilePath)
		if err != nil {
			fmt.Printf("Failed to open tile image %s: %v\n", tile.Path, err)
			continue
		}
		tileImgList = append(tileImgList, tileImg)
	}

	if len(tileImgList) == 0 {
		return errors.New("no valid tile images found")
	}

	// Update progress to 60%
	mosaic.Progress = 60
	db.DB.Save(mosaic)

	// Draw tiles on SD image
	for y := 0; y < sdHeight; y += tileSize / 2 {
		for x := 0; x < sdWidth; x += tileSize / 2 {
			// Get average color from main image for this region
			//avgColor := getAverageColor(mainImg, x*2, y*2, tileSize, tileSize)

			// Find best matching tile
			tileIdx := rand.Intn(len(tileImgList))
			tileImg := tileImgList[tileIdx]

			// Scale tile to fit
			scaledTile := scaleImage(tileImg, tileSize/2, tileSize/2)

			// Draw tile on SD image
			xdraw.Draw(sdImg, image.Rect(x, y, x+tileSize/2, y+tileSize/2), scaledTile, image.Point{0, 0}, draw.Over)
		}
	}

	// Update progress to 70%
	mosaic.Progress = 70
	db.DB.Save(mosaic)

	// Draw tiles on HD image
	for y := 0; y < hdHeight; y += tileSize {
		for x := 0; x < hdWidth; x += tileSize {
			// Get average color from main image for this region
			//avgColor := getAverageColor(mainImg, x, y, tileSize, tileSize)

			// Find best matching tile
			tileIdx := rand.Intn(len(tileImgList))
			tileImg := tileImgList[tileIdx]

			// Scale tile to fit
			scaledTile := scaleImage(tileImg, tileSize, tileSize)

			// Draw tile on HD image
			xdraw.Draw(hdImg, image.Rect(x, y, x+tileSize, y+tileSize), scaledTile, image.Point{0, 0}, draw.Over)
		}
	}

	// Update progress to 80%
	mosaic.Progress = 80
	db.DB.Save(mosaic)

	// Save the images
	if err := saveJPEG(sdImg, sdPath, 90); err != nil {
		return fmt.Errorf("failed to save SD image: %v", err)
	}

	// Update progress to 90%
	mosaic.Progress = 90
	db.DB.Save(mosaic)

	if err := saveJPEG(hdImg, hdPath, 90); err != nil {
		return fmt.Errorf("failed to save HD image: %v", err)
	}

	return nil
}

// GetMosaicStatus retrieves the status of a mosaic generation task
func (s *MosaicServiceImpl) GetMosaicStatus(userID uint, mosaicID uint) (*models.GeneratedMosaic, error) {
	var mosaic models.GeneratedMosaic
	result := db.DB.Where("id = ? AND user_id = ?", mosaicID, userID).First(&mosaic)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("mosaic not found")
		}
		return nil, result.Error
	}
	return &mosaic, nil
}

// GetProjectMosaics retrieves all mosaics for a project
func (s *MosaicServiceImpl) GetProjectMosaics(userID uint, projectID uint) ([]models.GeneratedMosaic, error) {
	var mosaics []models.GeneratedMosaic
	result := db.DB.Where("user_id = ? AND project_id = ?", userID, projectID).Order("created_at DESC").Find(&mosaics)
	if result.Error != nil {
		return nil, result.Error
	}
	return mosaics, nil
}

// Helper functions for image processing

// openImage opens an image file and returns an image.Image
func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// saveJPEG saves an image as JPEG with the specified quality
func saveJPEG(img image.Image, path string, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
}

// getAverageColor calculates the average color of a region in an image
func getAverageColor(img image.Image, x, y, width, height int) color.RGBA {
	var r, g, b, a uint32
	var count int

	bounds := img.Bounds()
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			px := x + dx
			py := y + dy
			if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
				pr, pg, pb, pa := img.At(px, py).RGBA()
				r += pr
				g += pg
				b += pb
				a += pa
				count++
			}
		}
	}

	if count == 0 {
		return color.RGBA{0, 0, 0, 255}
	}

	return color.RGBA{
		R: uint8(r / uint32(count) >> 8),
		G: uint8(g / uint32(count) >> 8),
		B: uint8(b / uint32(count) >> 8),
		A: uint8(a / uint32(count) >> 8),
	}
}

// scaleImage scales an image to the specified width and height
func scaleImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	xdraw.BiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}
