package controllers

import (
	"errors"
	"github.com/amityadav9314/goinkgrid/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	ee "github.com/pkg/errors"
)

func HandleHealthCheck(ctx *gin.Context) {
	curr := time.Now()
	type HealthCheckStatus struct {
		Status string `json:"status"`
	}

	var LOGGER = logger.GetLogger(ctx)
	health := HealthCheckStatus{Status: "UP"}
	//time.Sleep(3 * time.Second)
	LOGGER.Debug(ctx, "Debug test", "{'status': 'is_up'}", health, time.Since(curr))
	LOGGER.Info(ctx, "Health check up", "{'status': 'is_up'}", health, time.Since(curr))
	LOGGER.Error(ctx, "Some random error", "{'status': 'is_up'}", nil, time.Since(curr), ee.Wrap(errors.New("kuchbhi"), "error occcc"))
	LOGGER.Error(ctx, "Some random error", "{'status': 'is_up'}", nil, time.Since(curr), errors.New("kuchbhi"))
	ctx.JSON(http.StatusOK, health)
}
