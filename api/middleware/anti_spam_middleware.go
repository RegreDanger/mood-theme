package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var lastUpdate time.Time

func AntiSpamMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if time.Since(lastUpdate) < 1*time.Second {
			c.JSON(http.StatusOK, gin.H{"message": "Wait for the cooldown"})
			c.Abort()
		}
		lastUpdate = time.Now()
		c.Next()
	}
}
