package controllers

import (
	"gamesnight/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Need to send the right response code such as 404 if game not found or player not found
// Or if player is trying to get game that they are not part of
// Of if someone trying to make admin changes to game they are not admin of
func SendResponse(c *gin.Context, statusCode int, payload interface{}, err error) {
	resp := ApiResponse{
		Status:  "success",
		Message: "Request processed successfully",
		Data:    payload,
	}

	if err != nil {
		logger.GetLogger().Logger.Info("Response error", zap.Error(err))
		resp.Status = "error"
		resp.Message = "Request failed"
		resp.Error = err.Error()
	}

	c.Set("statusCode", statusCode)
	c.JSON(statusCode, resp)
}
