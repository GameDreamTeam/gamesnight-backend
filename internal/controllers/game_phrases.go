package controllers

import (
	"errors"
	"gamesnight/internal/models"
	"gamesnight/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPhrasesController(c *gin.Context) {
	player, err := getPlayerFromContext(c)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	gameId := c.Param("gameId")
	gameMeta, err := services.GetGameService().GetGameMeta(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	game, err := services.GetGameService().GetGame(gameId)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	if game.GameState != models.AddingWords {
		SendResponse(c, http.StatusConflict, nil, errors.New("The game is not in the AddingPhrasesState"))
		return
	}

	err = services.GetGameService().PlayerExistInGame(*gameMeta, *player)
	if err != nil {
		SendResponse(c, http.StatusNotFound, nil, err)
		return
	}

	err = services.GetPlayerService().PlayerAlreadyAddedPhrases(*player.Id)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	var phraseList models.PhraseList
	err = BindJSONAndHandleError(c, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	err = CheckPhraseListLength(phraseList)
	if err != nil {
		SendResponse(c, http.StatusBadRequest, nil, err)
		return
	}

	err = services.GetGameService().AddPhrasesToGame(*player.Id, gameMeta, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	err = services.GetGameService().AddPhrasesToPlayer(*player, &phraseList)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, nil, err)
		return
	}

	SendResponse(c, http.StatusOK, phraseList, nil)
}
