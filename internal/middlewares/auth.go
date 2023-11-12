package middlewares

import (
	"fmt"
	"gamesnight/internal/controllers"
	"gamesnight/internal/models"
	"gamesnight/internal/services"

	"github.com/gin-gonic/gin"
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
		fmt.Printf("Error in parsing token %s", err)
		return createPlayer(c)
	}
	return player, nil
}

// Moving token set to another function
func createPlayer(c *gin.Context) (*models.Player, error) {
	player, err := services.GetPlayerService().CreateNewPlayer()

	if err != nil {
		return nil, err
	}

	token, err := services.GetTokenService().CreatePlayerToken(*player.Id)
	if err != nil {
		return nil, err
	}

	// Need to check these other parameters
	c.SetCookie(playerCookieName, token.Token, 3600, "/", "", false, true)
	return player, nil
}
