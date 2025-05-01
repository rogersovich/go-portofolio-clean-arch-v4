package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

func Success(c *gin.Context, message string, data interface{}) {
	if data == nil {
		c.JSON(http.StatusOK, gin.H{
			"data":    []interface{}{}, // default empty array
			"message": message,
			"status":  "ok",
		})
		return
	}

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// If it's a slice and nil, return empty array
	if t.Kind() == reflect.Slice && v.IsNil() {
		c.JSON(http.StatusOK, gin.H{
			"data":    []interface{}{}, // empty JSON array
			"message": message,
			"status":  "ok",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"data":    nil,
		"message": message,
		"status":  "error",
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

type ReadingStats struct {
	WordCount        int     `json:"word_count"`
	TextLength       int     `json:"text_length"` // Number of characters (runes)
	EstimatedSeconds float64 `json:"estimated_seconds"`
	Minutes          int     `json:"minutes"` // Estimated minutes, rounded up
}

func RoundToOneDecimal(val float64) float64 {
	return math.Round(val*10) / 10
}

func ExtractHTMLtoStatistics(htmlContent string) ReadingStats {
	// WordsPerMinute is the average reading speed assumption.
	// Common values range from 200 to 250. Adjust as needed.
	const WordsPerMinute = 225.0

	// 1. Strip HTML tags to get plain text
	// Use bluemonday's StrictPolicy which removes all tags.
	p := bluemonday.StrictPolicy()
	plainText := p.Sanitize(htmlContent)

	// Trim leading/trailing whitespace for accurate word count
	plainText = strings.TrimSpace(plainText)

	// 2. Calculate Word Count
	// strings.Fields splits the string by whitespace into words
	words := strings.Fields(plainText)
	wordCount := len(words)

	// 3. Calculate Text Length (using RuneCount for multi-byte character safety)
	textLength := utf8.RuneCountInString(plainText)

	// 4. Calculate Estimated Reading Time
	var estimatedMinutesFloat float64
	if WordsPerMinute > 0 && wordCount > 0 {
		estimatedMinutesFloat = float64(wordCount) / WordsPerMinute
	} else {
		estimatedMinutesFloat = 0.0
	}

	estimatedSeconds := RoundToOneDecimal(estimatedMinutesFloat * 60.0)

	// Round minutes *up* to the nearest whole number (common practice for reading time)
	minutes := int(math.Ceil(estimatedMinutesFloat))

	// Handle the edge case where word count is 0 but we don't want 1 minute
	if wordCount == 0 {
		minutes = 0
	}

	// 5. Populate the result struct
	stats := ReadingStats{
		WordCount:        wordCount,
		TextLength:       textLength,
		EstimatedSeconds: estimatedSeconds,
		Minutes:          minutes,
	}

	return stats
}

func ConvertStringSliceToIntSlice(strs []string) ([]int, error) {
	ints := make([]int, len(strs))
	for i, s := range strs {
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, err // return error if any string is not a valid integer
		}
		ints[i] = num
	}
	return ints, nil
}

// ParseStringToTime is a helper function to convert a string to time.Time.
func ParseStringToTime(dateStr string, layout string) (time.Time, error) {
	// Parse the time string using the provided layout.
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %v", err)
	}
	return t, nil
}

func ParseStringPtrToTimePtr(dateStr *string, layout string) (*time.Time, error) {
	// If the string pointer is nil, return nil time pointer
	if dateStr == nil {
		return nil, nil
	}

	// Parse the time string using the provided layout
	t, err := time.Parse(layout, *dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}
	return &t, nil
}
