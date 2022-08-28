package middleware

import (
	"Newism/internal/model"
	"Newism/internal/utility"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utility.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse("Token Is Expire Or Incorrect", "Unauthorized"))
			c.Abort()
			return
		}
		c.Next()
	}
}
