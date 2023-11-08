package controllers

import (
	"fmt"
	"gamesnight/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewGameController(c *gin.Context) {

	database.GetRedis().Client.Set("abc", "asdas", time.Minute)

	fmt.Println(database.GetRedis().Client.Get("abcd"))
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from the controller!",
	})
}
