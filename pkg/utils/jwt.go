package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func getSecretJWT() []byte {
	envSecret := os.Getenv("JWT_SECRET")
	jwtSecret := []byte(envSecret)

	return jwtSecret
}

// GenerateJWT generates JWT token
func GenerateJWT(username string) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Set JWT secret
	jwtSecret := getSecretJWT()

	// Generate signed token string
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a JWT token
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// Set JWT secret
	jwtSecret := getSecretJWT()

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed with the correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
