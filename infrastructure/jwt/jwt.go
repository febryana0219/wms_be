package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var secretKey []byte

// Init loads the secret key from the environment variable
func Init() {
	_ = godotenv.Load() // Tidak perlu panic kalau .env tidak ditemukan (bisa di-load manual lewat env)
	secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		panic("JWT_SECRET_KEY environment variable not set or is empty")
	}
}

// GenerateToken generates a JWT token with userID and expiration time
func GenerateToken(userID string, expiration time.Time) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiration.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken validates the JWT token and returns the user ID if valid
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token claims")
	}

	return userID, nil
}

// GetUserIDFromContext extracts the user_id from Gin context (set by middleware)
func GetUserIDFromContext(c *gin.Context) (string, error) {
	val, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("user_id not found in context")
	}

	userID, ok := val.(string)
	if !ok {
		return "", errors.New("invalid user_id type in context")
	}

	return userID, nil
}
