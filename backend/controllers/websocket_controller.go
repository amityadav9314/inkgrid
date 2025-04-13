package controllers

import (
	"net/http"

	"github.com/amityadav9314/goinkgrid/logger"
	akyWs "github.com/amityadav9314/goinkgrid/pkg/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleWebSocketV2(ctx *gin.Context, pool *akyWs.Pool) {
	LOGGER := logger.GetLogger(ctx)
	conn, err := akyWs.Upgrade(ctx.Writer, ctx.Request)
	if err != nil {
		msg := "failed to upgrade: " + err.Error()
		LOGGER.Error(ctx, msg, nil, nil, 0, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	client := &akyWs.Client{
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client
	client.Read(ctx)
}

// Deprecated: HandleWebSocket
func HandleWebSocket(ctx *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	LOGGER := logger.GetLogger(ctx)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		msg := "failed to upgrade: " + err.Error()
		LOGGER.Error(ctx, msg, nil, nil, 0, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer conn.Close()

	for {
		messageType, message, err0 := conn.ReadMessage()
		LOGGER.Info(ctx, "Msg is read: "+string(message), nil, nil, 0)
		if err0 != nil {
			LOGGER.Error(ctx, "Error while reading msg", nil, nil, 0, err0)
			break
		}
		err1 := conn.WriteMessage(messageType, message)
		LOGGER.Info(ctx, "Msg is written: "+string(message), nil, nil, 0)
		if err1 != nil {
			LOGGER.Error(ctx, "Error while writing msg", nil, nil, 0, err1)
			break
		}
	}
}
