package routers

import (
	"gamesnight/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	r.GET("/health", controllers.HealthCheckController)

	api := r.Group("/v0")
	{
		// Game routes
		game := api.Group("/game")
		{
			game.GET("/create", controllers.NewGameController)
			game.POST("/:gameId/join", controllers.JoinGameController)
			game.GET("/:gameId/meta", controllers.GetGameMetaController)
			game.GET("/:gameId/details", controllers.GetGameController)
			game.POST("/:gameId/add-phrase", controllers.AddPhraseController)
			game.GET("/:gameId/phrases", controllers.GetGamePhrasesController)
			game.POST("/:gameId/make-teams", controllers.MakeTeamsController)
			game.GET("/:gameId/start", controllers.StartGameController)
			game.POST("/:gameId/start-turn", controllers.StartTurnController)
			game.POST("/:gameId/:playerId/remove", controllers.RemovePlayerController)
			// game.GET("/:gameId/play", controllers.GeneratePhraseController)
		}

		// Player routes
		player := api.Group("/player")
		{
			player.GET("/:playerId/phrases", controllers.GetPlayerPhrasesController)
		}

		// Feedback route
		api.POST("/feedback", controllers.SubmitFeedbackController)
	}
}
