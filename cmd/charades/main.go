package main

import (
	"gamesnight/internal/config"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/middlewares"
	"gamesnight/internal/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	//Setting up the environment
	config.New()

	//Initialization of zap logger
	logger.New()

	//Setting up Redis
	database.NewRedisClient()

	r := gin.New()

	//Setting up CORS, for frontend
	r.Use(middlewares.SetupCORS())

	r.Use(gin.Recovery())
	r.Use(middlewares.AuthMiddleware())
	r.Use(middlewares.ErrorHandlingMiddleware())
	r.Use(middlewares.LoggingMiddleware())
	routers.SetupRouter(r)

	r.Run(":8080")
}
