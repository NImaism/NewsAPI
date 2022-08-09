package Middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"time"
)

func RateLimit(Max float64, T int) gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(Max, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Minute * time.Duration(T)})
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}
