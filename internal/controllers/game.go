package controllers

import (
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGameController(c *gin.Context) {

	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}

	// Can check if this type conversion is passing or failing
	player := p.(*models.Player)
	game, err := services.GetGameService().CreateNewGame(player)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, game)
}

func JoinGameController(c *gin.Context) {
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	player := p.(*models.Player)

	gameId := c.Param("gameId")
	var playerName models.PlayerName

	if err := c.BindJSON(&playerName); err != nil {
		HandleError(c, err)
		return
	}

	player.Name = &playerName.Username

	game, err := services.GetGameService().JoinGame(gameId, player)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, game)
}

func GetGameController(c *gin.Context) {
	gameId := c.Param("gameId")
	game, err := services.GetGameService().GetGame(gameId)

	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, game)
}
