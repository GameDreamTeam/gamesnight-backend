package routers

import (
	"gamesnight/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	r.GET("/health", controllers.HealthCheckController)

	api := r.Group("/v0")
	{
		game := api.Group("/games")
		{
			game.POST("/", controllers.NewGameController)
			game.POST("/:gameId/join", controllers.JoinGameController)
			game.PATCH("/:gameId/update-state", controllers.UpdateState)

			game.GET("/:gameId/meta", controllers.GetGameMetaController)
			game.GET("/:gameId/details", controllers.GetGameController)
			game.GET("/:gameId/phrases", controllers.GetGamePhrasesController)
			game.GET("/:gameId/current-phrases", controllers.GetCurrentGamePhrasesController)

			game.POST("/:gameId/phrases", controllers.AddPhrasesController)
			game.POST("/:gameId/teams", controllers.MakeTeamsController)

			game.POST("/:gameId/start", controllers.StartGameController)
			game.POST("/:gameId/turns/start", controllers.StartTurnController)
			game.POST("/:gameId/choices", controllers.PlayerGuessController)
			game.POST("/:gameId/turns/end", controllers.EndTurnController)

			game.DELETE("/:gameId/players/:playerId", controllers.RemovePlayerController)

		}

		player := api.Group("/players")
		{
			//Maybe merge these 2 apis
			player.GET("/:playerId/phrases", controllers.GetPlayerPhrasesController)
			player.GET("/", controllers.GetPlayerDetailsController)
		}

		api.POST("/feedback", controllers.SubmitFeedbackController)
	}

	api2 := r.Group("/v1")
	{
		api2.POST("/feedback", controllers.SubmitFeedbackByEmailController)
	}
}
