package routers

import (
	"gamesnight/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	r.GET("/health", controllers.HealthCheckController)

	// Separate user and game routers
	r.GET("/v0/create-game", controllers.NewGameController)
	r.POST("/v0/:gameId/join", controllers.JoinGameController)
	r.GET("/v0/gamemeta/:gameId", controllers.GetGameMetaController)

	r.POST("/v0/game/:gameId/submit", controllers.AddPhraseController)
	r.GET("/v0/game/:gameId/phrases", controllers.GetGamePhrasesController)
	r.GET("/v0/player/:playerId/phrases", controllers.GetPlayerPhrasesController)

	// Maybe this should be a post call and not a get call
	r.GET("/v0/game/:gameId/start", controllers.StartGameController)
}
