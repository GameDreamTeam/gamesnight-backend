package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleError(c *gin.Context, err error) {

	//Move this to logger class
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Error("An error occurred",
		zap.Error(err),
	)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}
