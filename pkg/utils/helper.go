package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{
		"data":    data,
		"message": message,
		"status":  "ok",
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(201, gin.H{
		"data":    data,
		"message": message,
		"status":  "ok",
	})
}

func Error(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, gin.H{
		"data":    nil,
		"message": message,
		"status":  "error",
		"error":   err.Error(),
	})
}

func ErrorValidation(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
		"errors":  errors,
	})
}

func PaginatedSuccess(c *gin.Context, message string, data interface{}, page, limit, total int) {
	c.JSON(200, gin.H{
		"data": gin.H{
			"items": data,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
		"message": message,
		"status":  "ok",
		"error":   nil,
	})
}

func GetIsProduction() bool {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	return env == "production"
}

func GetProtocol() string {
	isProduction := GetIsProduction()
	if isProduction {
		return "https"
	}

	return "http" // default development
}

func PrintJSON(v any) {
	jsonBytes, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(jsonBytes))
}
