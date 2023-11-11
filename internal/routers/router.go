package routers

import (
	"gamesnight/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	r.GET("/health", controllers.HealthCheckController)

	r.GET("/v0/create-game", controllers.NewGameController)
	r.POST("/v0/:gameId/join", controllers.JoinGameController)
	r.GET("/v0/:gameId", controllers.GetGameController)

}
