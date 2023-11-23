package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SendMessageToSlack(c *gin.Context) {
	var message struct {
		Text string `json:"text"`
	}

	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Use the Slack webhook URL you obtained in Task 1
	webhookURL := "https://hooks.slack.com/services/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"

	err := sendToSlack(webhookURL, message.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message to Slack"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Message sent to Slack successfully"})
}

func sendToSlack(webhookURL, message string) error {
	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(fmt.Sprintf(`{"text": "%s"}`, message)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Slack API Response:", string(body))

	return nil
}
