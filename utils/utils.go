package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"productanalyzer/api/config"
	api_error "productanalyzer/api/errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func FormatValidationErrors(err error) interface{} {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return gin.H{"error": "Invalid input."}
	}
	errorMap := map[string][]string{
		"required": {},
		"email":    {},
		"invalid":  {},
	}
	for _, fieldError := range validationErrors {
		fieldName := fieldError.Field()
		switch fieldError.Tag() {
		case "required":
			errorMap["required"] = append(errorMap["required"], fieldName)
		case "email":
			errorMap["email"] = append(errorMap["email"], fieldName)
		default:
			errorMap["invalid"] = append(errorMap["invalid"], fieldName)
		}
	}

	return errorMap
}

func CreateToken(userID string) (string, *api_error.APIError) {
	expirationTime := time.Now().UTC().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sign, err := token.SignedString([]byte(config.Config.SECRET_KEY))
	if err != nil {
		return "", api_error.UnexpectedError(err)
	}
	return sign, nil
}

func ValidateToken(tokenString string) (*Claims, *api_error.APIError) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.SECRET_KEY), nil
	})

	if err != nil {
		return nil, api_error.UnexpectedError(err)
	}

	if !token.Valid {
		return nil, api_error.NewAPIError("Login Failed!", http.StatusUnauthorized, "Invalid or expired token")
	}

	return claims, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword checks if the provided password matches the hashed password
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Generate numeric one-time password of given length
func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}
