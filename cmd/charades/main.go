package main

import (
	"gamesnight/internal/controllers"
	"gamesnight/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.NewRedisClient()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/v0/create-game", controllers.NewGameController)

	router.Run(":8080")
}
