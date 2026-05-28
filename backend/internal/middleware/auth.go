package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_id") == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "не авторизован"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_id") == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "не авторизован"})
			c.Abort()
			return
		}
		role, _ := session.Get("role").(string)
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "доступ запрещён"})
			c.Abort()
			return
		}
		c.Next()
	}
}
