package utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		duration := time.Since(start)
		Logger.WithFields(map[string]interface{}{
			"method":   method,
			"path":     path,
			"status":   c.Writer.Status(),
			"duration": duration,
		}).Info("request-log")
	}
}

func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				Logger.WithFields(logrus.Fields{
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
					"error":  rec,
				}).Error("🔥 Panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"data":    nil,
					"message": "Internal server error",
					"status":  "error",
					"error":   rec,
				})
			}
		}()
		c.Next()
	}
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			Error(c, http.StatusBadRequest, "Unauthorized request")
			c.Abort()
			return
		}

		// Remove "Bearer " from the token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validate the token
		token, err := ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			Error(c, http.StatusBadRequest, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set the token claims in the context for further use if needed (e.g., username)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
			c.Set("email", claims["email"])
		}

		// Continue to the next handler
		c.Next()
	}
}
