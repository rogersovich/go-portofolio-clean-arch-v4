package utils

import (
	"net/http"
	"time"

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
