package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	akyWs "github.com/amityadav9314/goinkgrid/pkg/websocket"

	"github.com/amityadav9314/goinkgrid/config"
	"github.com/amityadav9314/goinkgrid/constants"
	"github.com/amityadav9314/goinkgrid/logger"
	"github.com/amityadav9314/goinkgrid/routers"
	"github.com/amityadav9314/goinkgrid/utils"
	"github.com/gin-contrib/cors"
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

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"} // Ensure the correct origin is included
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "OPTIONS"}
	corsConfig.ExposeHeaders = []string{"*"} // Expose all headers
	corsConfig.AllowCredentials = true       // Allow credentials if needed
	mainRouter.Use(cors.New(corsConfig))

	// Handle OPTIONS method for preflight requests
	mainRouter.OPTIONS("/goinkgrid/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Initialize websocket pool
	var pool = akyWs.NewPool()
	go pool.Start()

	routers.InitRoutes(mainRouter, ENVIRONMENT, pool)

	log.Println("Server is running on port:", PORT)
	err := mainRouter.Run(":" + PORT)
	if err != nil {
		log.Println("Error in starting server: ", err)
		return
	}
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
