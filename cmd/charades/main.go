package main

import (
	"gamesnight/internal/config"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/middlewares"
	"gamesnight/internal/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.New()
	logger.New()
	database.NewRedisClient()

	r := gin.New()

	// Move this logic to another file
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	r.Use(gin.Recovery())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.AuthMiddleware())

	routers.SetupRouter(r)

	r.Run(":8080")
}
