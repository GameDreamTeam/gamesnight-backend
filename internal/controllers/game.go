package controllers

import (
	"gamesnight/internal/database"
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

func AddPhraseController(c *gin.Context) {
	var phraseList models.PhraseList
	gameId := c.Param("gameId")

	// Assuming the player ID is available in the context
	p, exists := c.Get("player")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	player := p.(*models.Player)
	playerId := *player.Id

	// Bind the incoming JSON to phraseList
	if err := c.BindJSON(&phraseList); err != nil {
		HandleError(c, err)
		return
	}

	// Save the phrases in Redis with gameId as the key
	err := database.SetGamePhrases(gameId, &phraseList)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Optionally, save the phrases with a key combining the gameId and playerId
	err = database.SetPlayerGamePhrases(gameId, playerId, &phraseList)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Phrases added successfully"})
}

func GetGamePhrasesController(c *gin.Context) {
	gameId := c.Param("gameId")
	phrases, err := database.GetGamePhrases(gameId)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, phrases)
}

func GetPlayerPhrasesController(c *gin.Context) {
	playerId := c.Param("playerId")
	gameId := c.Query("gameId") // Assuming gameId is also required to fetch specific phrases

	phrases, err := database.GetPlayerGamePhrases(gameId, playerId)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, phrases)
}
