package middlewares

import (
	"gamesnight/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		logger.GetLogger().Logger.Info(
			"Request Logger",
			zap.String("Type", "http-request"),
			zap.String("Method", c.Request.Method),
			zap.String("Path", c.Request.URL.Path),
			zap.Duration("Duration", latency),
		)
	}
}
