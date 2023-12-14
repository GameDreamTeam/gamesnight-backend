package middlewares

import (
	"fmt"
	"gamesnight/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", r)
				}
				controllers.SendResponse(c, http.StatusInternalServerError, nil, err)
				c.Abort()
			}
		}()
		c.Next()
	}
}
