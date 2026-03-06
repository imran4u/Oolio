package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		key := c.GetHeader("api_key")

		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "missing api key",
			})
			return
		}

		if key != "apitest" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid api key",
			})
			return
		}

		c.Next()
	}
}
