package app

import (
	"github.com/amityadav9314/goinkgrid/internal/api/handlers"
	"github.com/amityadav9314/goinkgrid/internal/services"
	akyWs "github.com/amityadav9314/goinkgrid/pkg/websocket"
	"gorm.io/gorm"
)

// ServiceProvider centrally manages all application services and handlers
type ServiceProvider struct {
	db        *gorm.DB
	jwtSecret string
	pool      *akyWs.Pool

	// Services
	userService    services.UserService
	projectService services.ProjectService
	imageService   services.ImageService
	mosaicService  services.MosaicService

	// Handlers
	authHandler    *handlers.AuthHandler
	projectHandler *handlers.ProjectHandler
	imageHandler   *handlers.ImageHandler
	mosaicHandler  *handlers.MosaicHandler
}

// NewServiceProvider initializes the service provider with dependencies
func NewServiceProvider(db *gorm.DB, jwtSecret string, pool *akyWs.Pool) *ServiceProvider {
	sp := &ServiceProvider{
		db:        db,
		jwtSecret: jwtSecret,
		pool:      pool,
	}

	// Initialize services
	sp.initServices()

	// Initialize handlers
	sp.initHandlers()

	return sp
}

// initServices initializes all services
func (sp *ServiceProvider) initServices() {
	sp.userService = services.NewUserService(sp.db)
	sp.projectService = services.NewProjectService(sp.db)
	sp.imageService = services.NewImageService(sp.db)
	sp.mosaicService = services.NewMosaicService("./uploads")
}

// initHandlers initializes all handlers
func (sp *ServiceProvider) initHandlers() {
	sp.authHandler = handlers.NewAuthHandler(sp.userService, sp.jwtSecret)
	sp.projectHandler = handlers.NewProjectHandler(sp.projectService)
	sp.imageHandler = handlers.NewImageHandler("./uploads", sp.imageService)
	sp.mosaicHandler = handlers.NewMosaicHandler(sp.mosaicService)
}

// UserService returns the user service
func (sp *ServiceProvider) UserService() services.UserService {
	return sp.userService
}

// ProjectService returns the project service
func (sp *ServiceProvider) ProjectService() services.ProjectService {
	return sp.projectService
}

// ImageService returns the image service
func (sp *ServiceProvider) ImageService() services.ImageService {
	return sp.imageService
}

// MosaicService returns the mosaic service
func (sp *ServiceProvider) MosaicService() services.MosaicService {
	return sp.mosaicService
}

// AuthHandler returns the auth handler
func (sp *ServiceProvider) AuthHandler() *handlers.AuthHandler {
	return sp.authHandler
}

// ProjectHandler returns the project handler
func (sp *ServiceProvider) ProjectHandler() *handlers.ProjectHandler {
	return sp.projectHandler
}

// ImageHandler returns the image handler
func (sp *ServiceProvider) ImageHandler() *handlers.ImageHandler {
	return sp.imageHandler
}

// MosaicHandler returns the mosaic handler
func (sp *ServiceProvider) MosaicHandler() *handlers.MosaicHandler {
	return sp.mosaicHandler
}

// JWTSecret returns the JWT secret
func (sp *ServiceProvider) JWTSecret() string {
	return sp.jwtSecret
}

// Pool returns the websocket pool
func (sp *ServiceProvider) Pool() *akyWs.Pool {
	return sp.pool
}
