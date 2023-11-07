package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware: logger and recovery middleware.
	router := gin.Default()

	// Define a simple GET route.
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start the server on port 8080.
	router.Run(":8080")
}
