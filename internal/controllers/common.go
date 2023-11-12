package controllers

import (
	"gamesnight/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleError(c *gin.Context, err error) {

	logger.GetLogger().Logger.Error("Server error", zap.Error(err))
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}
