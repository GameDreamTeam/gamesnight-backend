package controllers

import (
	"errors"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SubmitFeedbackByEmailController(c *gin.Context) {

	var message models.Message

	if err := c.BindJSON(&message); err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	// Instantiate EmailService
	emailService := services.GetEmailService()

	// Call the method on the EmailService instance
	err := emailService.SendEmail("parusgiri@gmail.com", "Feedback for GamesNight", message.Text)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, errors.New("failed to send email"))
		return
	}

	SendResponse(c, http.StatusOK, "Email sent successfully", nil)
}
