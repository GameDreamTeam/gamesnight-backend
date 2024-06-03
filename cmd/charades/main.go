package main

import (
	"gamesnight/internal/config"
	"gamesnight/internal/database"
	"gamesnight/internal/logger"
	"gamesnight/internal/middlewares"
	"gamesnight/internal/routers"
	"gamesnight/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	//Setting up the environment
	config.New()
	godotenv.Load()

	//Initialization of zap logger
	logger.New()

	//Setting up Redis
	database.NewRedisClient()

	//Email-Service
	email := "ayushgupta71011@gmail.com"
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	services.NewEmailService(email, password, smtpHost, smtpPort)

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
