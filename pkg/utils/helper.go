package utils

import (
	"encoding/json"
	"errors"
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

func BoolToYN(val bool) string {
	if val {
		return "Y"
	}
	return "N"
}

func StringBoolToYN(val string) string {
	if val == "1" {
		return "Y"
	}
	return "N"
}

// Validates that a string array is not empty and doesn't contain only empty strings
func ValidateFormArrayString(strs []string, field string, is_required bool) ([]string, error) {
	var result []string

	for _, s := range strs {
		if s != "" {
			result = append(result, s)
		}
	}

	if is_required && len(result) == 0 {
		return nil, errors.New(field + " array must not be empty")
	}

	// Ensure empty slice (not nil) is returned if not required
	if !is_required && len(result) == 0 {
		return []string{}, nil
	}

	return result, nil
}

// BuildSQLInClause generates a string of "?, ?, ?" placeholders and a slice of interface{} args
func BuildSQLInClause[intType ~int | ~int64 | ~string](values []intType) (string, []interface{}) {
	placeholders := make([]string, len(values))
	args := make([]interface{}, len(values))

	for i, v := range values {
		placeholders[i] = "?"
		args[i] = v
	}

	return strings.Join(placeholders, ","), args
}

func ToInterfaceSlice(arr []interface{}) []interface{} {
	out := make([]interface{}, len(arr))
	for i, v := range arr {
		out[i] = v
	}
	return out
}
