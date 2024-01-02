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

			game.GET("/:gameId/meta", controllers.GetGameMetaController)
			game.GET("/:gameId/details", controllers.GetGameController)
			game.GET("/:gameId/phrases", controllers.GetGamePhrasesController)

			game.POST("/:gameId/phrases", controllers.AddPhraseController)
			game.POST("/:gameId/teams", controllers.MakeTeamsController)
			game.POST("/:gameId/start", controllers.StartGameController)
			game.POST("/:gameId/turns/start", controllers.StartTurnController)

			game.DELETE("/:gameId/players/:playerId", controllers.RemovePlayerController)

			game.POST("/:gameId/choices", controllers.PlayerGuessController)
			game.POST("/:gameId/turns/end", controllers.EndTurnController)
		}

		player := api.Group("/players")
		{
			//Maybe merge these 2 apis
			player.GET("/:playerId/phrases", controllers.GetPlayerPhrasesController)
			player.GET("/:playerId", controllers.GetPlayerDetailsController)
		}

		api.POST("/feedback", controllers.SubmitFeedbackController)
	}
}
