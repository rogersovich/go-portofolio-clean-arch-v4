package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []FieldError `json:"errors"`
}

var customMessages = map[string]string{
	"name.required":        "name is required",
	"title.required":       "title is required",
	"avatar_file.required": "avatar_file is required",
	"description.required": "description is required",
	"is_major.oneof":       "is_major must be either 'Y' or 'N'",
	"id.required":          "id is required",
	"id.numeric":           "id must be numeric",
}

// ValidateStruct handles JSON binding and validation errors, and returns a structured error response.
// Returns true if validation passed, or false if errors are returned to the client.
func ValidateStruct(c *gin.Context, requestStruct interface{}, bindErr error) bool {
	if bindErr == nil {
		return true // No validation error, continue execution
	}

	var validationErrors validator.ValidationErrors

	// Check if the error is a validation error
	if errors.As(bindErr, &validationErrors) {
		var formattedErrors []FieldError

		for _, fieldError := range validationErrors {
			// Default field name (fallback to struct field name)
			jsonField := fieldError.Field()

			// Try to get the actual JSON tag from the struct
			if structField, ok := reflect.TypeOf(requestStruct).Elem().FieldByName(fieldError.StructField()); ok {
				jsonTag := structField.Tag.Get("json")
				if jsonTag != "" {
					jsonField = jsonTag
				}
			}

			// Create the key for custom message lookup (e.g. "tech_name.required")
			messageKey := fmt.Sprintf("%s.%s", jsonField, fieldError.Tag())

			// Look up the custom error message or fallback to default message
			errorMessage, exists := customMessages[messageKey]
			if !exists {
				if fieldError.Tag() == "oneof" {
					params := formatTypeOf(fieldError.Param())
					errorMessage = fmt.Sprintf("%s must be one of '%s'", jsonField, params)
				} else {
					errorMessage = fmt.Sprintf("%s failed on '%s'", jsonField, fieldError.Tag())
				}
			}

			// Append the error in the desired format
			formattedErrors = append(formattedErrors, FieldError{
				Field:   jsonField,
				Message: errorMessage,
			})
		}

		// Return all validation errors in the expected JSON format
		c.JSON(400, gin.H{
			"errors": formattedErrors,
		})
		return false
	}

	// Fallback for non-validation binding errors (e.g. malformed JSON)
	switch err := bindErr.(type) {
	case *json.SyntaxError:
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Malformed JSON at offset %d", err.Offset),
		})
		return false
	case *json.UnmarshalTypeError:
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Wrong type for field '%s': %s", err.Field, err.Error()),
		})
		return false
	default:
		c.JSON(400, gin.H{
			"error": "Invalid request body: " + bindErr.Error(),
		})
		return false
	}
}

func ValidateRequest(data interface{}) []FieldError {
	var validate = validator.New()

	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var errs []FieldError
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	for _, fe := range err.(validator.ValidationErrors) {
		fieldName := fe.StructField()
		tag := fe.Tag()

		// get json tag
		field, _ := typ.FieldByName(fieldName)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = strings.ToLower(fieldName)
		} else {
			jsonTag = strings.Split(jsonTag, ",")[0] // handle omitempty
		}

		// get message from custom map
		customKey := fmt.Sprintf("%s.%s", jsonTag, tag)
		message, ok := customMessages[customKey]
		if !ok {
			// fallback message
			message = fmt.Sprintf("%s is not valid (%s)", jsonTag, tag)
		}

		errs = append(errs, FieldError{
			Field:   jsonTag,
			Message: message,
		})
	}

	return errs
}

func formatTypeOf(typeof string) string {
	parts := strings.Fields(typeof)  // splits by any whitespace
	return strings.Join(parts, ", ") // joins with ", "
}
