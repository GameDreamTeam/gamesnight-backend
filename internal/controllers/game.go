package controllers

import (
	"fmt"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

const userCookieName = "sid1"

func NewGameController(c *gin.Context) {
	userCookie, err := c.Cookie(userCookieName)

	if err != nil {
		user, err := services.GetUserService().CreateNewUser()

		if err != nil {
			fmt.Printf("Error in creating new user %s", err)
			// We have to remove these errors from UI when exposing to users
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}

		token, err := services.GetTokenService().CreateUserToken(*user.UserId)
		if err != nil {
			fmt.Printf("Error in creating new token %s", err)
			// We have to remove these errors from UI when exposing to users
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}
		// Need to check these other parameters
		c.SetCookie(userCookieName, token.Token, 3600, "/", "", false, true)
		userCookie = token.Token
	}

	user, err := services.GetTokenService().ParseUserToken(userCookie)
	if err != nil {
		fmt.Printf("Error in parsing token %s", err)
		// We have to remove these errors from UI when exposing to users
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	game, err := services.GetGameService().CreateNewGame(user)
	if err != nil {
		fmt.Printf("Error in creating game %s", err)
		// We have to remove these errors from UI when exposing to users
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, game)
}
