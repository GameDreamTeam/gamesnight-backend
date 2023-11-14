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
	r.GET("/v0/game/:gameId", controllers.GetGameController)

	// Maybe this should be a post call and not a get call
	// Maybe using hyphens is not good practice
	r.GET("/v0/game/:gameId/divide-teams", controllers.MakeTeamsController)
	r.GET("/v0/game/:gameId/start", controllers.StartGameController)

}
