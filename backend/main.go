package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/amityadav9314/goinkgrid/config"
	"github.com/amityadav9314/goinkgrid/constants"
	"github.com/amityadav9314/goinkgrid/internal/app"
	db "github.com/amityadav9314/goinkgrid/internal/db/postgres"
	"github.com/amityadav9314/goinkgrid/logger"
	akyWs "github.com/amityadav9314/goinkgrid/pkg/websocket"
	"github.com/amityadav9314/goinkgrid/routers"
	"github.com/amityadav9314/goinkgrid/utils"
	"github.com/gin-gonic/gin"
)

var (
	ENVIRONMENT string
	PORT        string
)

func main() {
	initialize()
	gin.SetMode(gin.DebugMode)
	mainRouter := gin.Default()
	mainRouter.MaxMultipartMemory = 10 << 20 // 100 MiB

	// Configure CORS - Use the most permissive configuration for development
	mainRouter.Use(func(c *gin.Context) {
		// Allow all origins in development
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		// Set CORS headers with maximum permissiveness for development
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Global OPTIONS handler for preflight requests
	mainRouter.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Initialize database
	db.Init()

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-jwt-secret-key" // For development only
	}

	// Initialize websocket pool
	pool := akyWs.NewPool()
	go pool.Start()

	// Initialize service provider with dependencies
	serviceProvider := app.NewServiceProvider(db.DB, jwtSecret, pool)

	// Set up routes with the service provider
	routers.InitRoutes(mainRouter, ENVIRONMENT, serviceProvider)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mainRouter,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server is running on port:", PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func initialize() {
	setEnvironment()
	config.DoInit(ENVIRONMENT)
	logger.InitLogger()
}

func setEnvironment() {
	environmentList := []string{constants.EnvDev, constants.EnvPp, constants.EnvProdpp, constants.EnvProd}
	envFlagPtr := flag.String("env", "dev", "environment: "+(strings.Join(environmentList, ",")))
	portFlagPtr := flag.String("port", "8034", "8034")

	flag.Parse()
	PORT = *portFlagPtr
	if utils.StringInSlice(*envFlagPtr, environmentList) {
		config.SetEnv(*envFlagPtr)
		ENVIRONMENT = *envFlagPtr
	} else {
		config.SetEnv(constants.EnvDev)
		ENVIRONMENT = constants.EnvDev
	}
	log.Println("ENVIRONMENT: " + config.GetEnv())
}
