package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"

	"net/http"
)

func SubmitFeedback(c *gin.Context) {
	// Assuming feedback is sent as JSON in the request body
	var feedback struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Send the feedback message to Slack
	err := sendFeedbackToSlack(feedback.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Feedback submitted successfully"})
}
func sendFeedbackToSlack(message string) error {
	// Use the Slack webhook URL you obtained in Task 2
	webhookURL := "https://hooks.slack.com/services/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"

	api := slack.New(
		// You can use an OAuth token instead of a webhook URL for more advanced features
		webhookURL,
	)

	attachment := slack.Attachment{
		Text:  message,
		Color: "#36a64f",
	}

	channelID, timestamp, err := api.PostMessage(webhookURL,
		slack.MsgOptionText("New Feedback", false),
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	return nil
}
