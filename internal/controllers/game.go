package controllers

import (
	"fmt"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

const userCookieName = "sid1"

func NewGameController(c *gin.Context) {
	user, err := services.GetUserService().GetUser(c, userCookieName)
	if err != nil {
		// Move this to error handler middle ware
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	// We have to check whether user has an existing game or not
	// If he has then we should allow him to select that instead of creating new game
	game, err := services.GetGameService().CreateNewGame(user)
	if err != nil {
		// Move this to error handler middle ware
		fmt.Printf("Error in creating game %s", err)
		// We have to remove these errors from UI when exposing to users
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, game)
}

func JoinGameController(c *gin.Context) {
	user, err := services.GetUserService().GetUser(c, userCookieName)

	if err != nil {
		fmt.Printf("Error in Getting user %s", err)
		// We have to remove these errors from UI when exposing to users
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	gameId := c.Param("gameId")
	fmt.Println(gameId)
	game, err := services.GetGameService().JoinGame(gameId, *user.UserId)
	if err != nil {
		// Move this to error handler middle ware
		fmt.Printf("Error in joining game %s", err)
		// We have to remove these errors from UI when exposing to users
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, game)
}
