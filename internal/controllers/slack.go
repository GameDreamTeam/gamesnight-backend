package controllers

import (
	"errors"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

const webhookURL = "https://hooks.slack.com/services/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"

func SubmitFeedbackController(c *gin.Context) {

	var message models.Message

	if err := c.BindJSON(&message); err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	// Instantiate SlackService
	slackService := services.GetSlackService()

	// Call the method on the SlackService instance
	err := slackService.SendToSlack(webhookURL, message.Text)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, errors.New("Failed to send message to Slack"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Message sent successfully"})
}
