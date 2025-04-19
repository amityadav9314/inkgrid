package routers

import (
	"github.com/amityadav9314/goinkgrid/controllers"
	"github.com/amityadav9314/goinkgrid/internal/api/handlers"
	"github.com/amityadav9314/goinkgrid/internal/api/middleware"
	akyWs "github.com/amityadav9314/goinkgrid/pkg/websocket"
	"github.com/gin-gonic/gin"
)

func InitRoutes(mainRouter *gin.Engine, environment string, pool *akyWs.Pool) {
	// Serve static files from uploads directory
	mainRouter.Static("/uploads", "./uploads")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler("your-jwt-secret-key")
	imageHandler := handlers.NewImageHandler("./uploads")
	mosaicHandler := handlers.NewMosaicHandler()
	projectHandler := handlers.NewProjectHandler()

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware("your-jwt-secret-key")

	// Base API group
	api := mainRouter.Group("/goinkgrid")

	// Existing routes
	api.GET("/health", controllers.HandleHealthCheck)
	api.GET("/ws", controllers.HandleWebSocket)
	api.GET("/v2/ws", func(c *gin.Context) {
		controllers.HandleWebSocketV2(c, pool)
	})

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
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
				imagesAuth.POST("/main", imageHandler.UploadMainImage)
				imagesAuth.POST("/tiles", imageHandler.UploadTileImages)
				imagesAuth.GET("/tiles", imageHandler.GetTileCollections)
			}
		}

		// Project routes (all require auth)
		projects := apiV1.Group("/projects") // Removed trailing slash to prevent redirects
		projects.Use(authMiddleware.RequireAuth())
		{
			projects.GET("/", projectHandler.ListProjects)
			projects.POST("/", projectHandler.CreateProject)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
		}

		// Mosaic generation routes (all require auth)
		generate := apiV1.Group("/generate")
		generate.Use(authMiddleware.RequireAuth())
		{
			generate.POST("/", mosaicHandler.GenerateMosaic)
			generate.GET("/:id/status", mosaicHandler.GetGenerationStatus)
		}
	}
}
