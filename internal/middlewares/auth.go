package middlewares

import (
	"gamesnight/internal/controllers"
	"gamesnight/internal/logger"
	"gamesnight/internal/models"
	"gamesnight/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const playerCookieName = "sid1"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		p, err := getPlayerOrCreateNew(c)
		if err != nil {
			controllers.HandleError(c, err)
			return
		}

		c.Set("player", p)
		c.Next()
	}
}

func getPlayerOrCreateNew(c *gin.Context) (*models.Player, error) {
	playerCookie, err := c.Cookie(playerCookieName)
	if err != nil {
		return createPlayer(c)
	}

	player, err := services.GetTokenService().ParsePlayerToken(playerCookie)
	if err != nil {
		logger.GetLogger().Logger.Error("Auth token", zap.Error(err))
		return createPlayer(c)
	}
	return player, nil
}

func createPlayer(c *gin.Context) (*models.Player, error) {
	player, err := services.GetPlayerService().CreateNewPlayer()

	if err != nil {
		return nil, err
	}

	token, err := services.GetTokenService().CreatePlayerToken(*player.Id)
	setUserAuthCookie(c, token)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func setUserAuthCookie(c *gin.Context, token *models.Token) {
	c.SetCookie(playerCookieName, token.Token, 3600, "/", "", false, true)
}
