package routers

import (
	"github.com/amityadav9314/goinkgrid/controllers"
	"github.com/amityadav9314/goinkgrid/internal/api/middleware"
	"github.com/amityadav9314/goinkgrid/internal/app"
	"github.com/gin-gonic/gin"
)

func InitRoutes(mainRouter *gin.Engine, environment string, serviceProvider *app.ServiceProvider) {
	// Serve static files from uploads directory
	mainRouter.Static("/uploads", "./uploads")

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(serviceProvider.JWTSecret())

	// Base API group
	api := mainRouter.Group("/goinkgrid")

	// Existing routes
	api.GET("/health", controllers.HandleHealthCheck)
	api.GET("/ws", controllers.HandleWebSocket)
	api.GET("/v2/ws", func(c *gin.Context) {
		controllers.HandleWebSocketV2(c, serviceProvider.Pool())
	})

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", serviceProvider.AuthHandler().Register)
		auth.POST("/login", serviceProvider.AuthHandler().Login)
		auth.POST("/refresh", serviceProvider.AuthHandler().RefreshToken)
	}

	// API routes that require authentication
	apiV1 := api.Group("/api")
	{
		// Public image routes
		images := apiV1.Group("/images")
		{
			// Protected routes
			imagesAuth := images.Group("/")
			imagesAuth.Use(authMiddleware.RequireAuth())
			{
				imagesAuth.POST("/main", serviceProvider.ImageHandler().UploadMainImage)
				imagesAuth.POST("/tiles", serviceProvider.ImageHandler().UploadTileImages)
				imagesAuth.GET("/tiles", serviceProvider.ImageHandler().GetTileCollections)
			}
		}

		// Project routes (all require auth)
		projects := apiV1.Group("/projects")
		projects.Use(authMiddleware.RequireAuth())
		{
			projects.GET("/", serviceProvider.ProjectHandler().ListProjects)
			projects.POST("/", serviceProvider.ProjectHandler().CreateProject)
			projects.GET("/:id", serviceProvider.ProjectHandler().GetProject)
			projects.PUT("/:id", serviceProvider.ProjectHandler().UpdateProject)
			projects.DELETE("/:id", serviceProvider.ProjectHandler().DeleteProject)
		}

		// Mosaic generation routes (all require auth)
		generate := apiV1.Group("/generate")
		generate.Use(authMiddleware.RequireAuth())
		{
			generate.POST("/", serviceProvider.MosaicHandler().GenerateMosaic)
			generate.GET("/:id/status", serviceProvider.MosaicHandler().GetGenerationStatus)
			generate.POST("/settings", serviceProvider.MosaicHandler().SaveMosaicSettings)
			generate.GET("/settings", serviceProvider.MosaicHandler().GetMosaicSettings)
		}
	}
}
