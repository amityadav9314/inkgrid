package routers

import (
	"github.com/amityadav9314/goinkgrid/controllers"
	akyWs "github.com/amityadav9314/goinkgrid/pkg/websocket"
	"github.com/gin-gonic/gin"
)

func InitRoutes(mainRouter *gin.Engine, environment string, pool *akyWs.Pool) {
	api := mainRouter.Group("/goinkgrid")
	api.GET("/health", controllers.HandleHealthCheck)
	api.GET("/ws", controllers.HandleWebSocket)
	api.GET("/v2/ws", func(c *gin.Context) {
		controllers.HandleWebSocketV2(c, pool)
	})
}
