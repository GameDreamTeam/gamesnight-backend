package controllers

import (
	"errors"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getPlayerFromContext(c *gin.Context) (*models.Player, error) {
	p, exists := c.Get("player")
	if !exists {
		return nil, errors.New("player does not exist")
	}
	return p.(*models.Player), nil
}

func isAdminPlayer(gamemeta models.GameMeta, player *models.Player) error {
	if *player.Id != gamemeta.AdminId {
		logger.GetLogger().Logger.Error(
			"player starting game should be admin",
			zap.Any("gamemeta", gamemeta),
			zap.Any("player", player),
		)
		return errors.New("player starting game should be admin")
	}

	return nil
}

func BindJSONAndHandleError(c *gin.Context, obj interface{}) error {
	if err := c.BindJSON(obj); err != nil {
		return errors.New("cannot bind body data")
	}
	return nil
}

func CheckPhraseListLength(phraseList models.PhraseList) error {
	if len(*phraseList.List) != 4 {
		return errors.New("number of phrases must be four")
	}
	return nil
}
