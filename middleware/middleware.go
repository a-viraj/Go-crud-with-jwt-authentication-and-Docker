package middleware

import (
	"net/http"

	"github.com/a-viraj/project/helper"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token string missing"})
			c.Abort()
		}
		claims, err := helper.ValidateTokens(token)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("userid", claims.UserId)
		c.Set("usertype", claims.UserType)
		c.Next()
	}
}
