package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"productanalyzer/api/config"
	api_error "productanalyzer/api/errors"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	DEVICE_TYPE_MOBILE  = "mobile"
	DEVICE_TYPE_TABLET  = "tablet"
	DEVICE_TYPE_DESKTOP = "desktop"
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

// GetDeviceType returns the type of device based on the user agent
func GetDeviceType(r *http.Request) string {
	userAgent := strings.ToLower(r.UserAgent())

	mobilePatterns := []string{
		"android", "webos", "iphone", "ipod", "blackberry", "windows phone",
		"opera mini", "iemobile", "mobile", "nokia", "fennec", "kindle", "silk",
		"maemo", "palm", "midp", "samsung", "symbian", "j2me", "wap", "avantgo",
		"blazer", "plucker", "xiino", "vodafone", "docomo", "softbank",
	}

	tabletPatterns := []string{
		"ipad", "tablet", "playbook", "silk", "kindle", "puffin",
	}

	for _, pattern := range mobilePatterns {
		if strings.Contains(userAgent, pattern) {
			if pattern == "iphone" && strings.Contains(userAgent, "ipad") {
				continue
			}
			return DEVICE_TYPE_MOBILE
		}
	}

	for _, pattern := range tabletPatterns {
		if strings.Contains(userAgent, pattern) {
			return DEVICE_TYPE_TABLET
		}
	}

	tabletUserAgents := []string{
		"(?i)android.*nexus\\s*7",
		"(?i)android.*nexus\\s*9",
		"(?i)android.*nexus\\s*10",
		"(?i)android.*sm-t\\d{3}",
		"(?i)android.*gt-p\\d{4}",
		"(?i)android.*sch-i\\d{3}",
	}

	for _, pattern := range tabletUserAgents {
		match, _ := regexp.MatchString(pattern, userAgent)
		if match {
			return DEVICE_TYPE_TABLET
		}
	}
	return DEVICE_TYPE_DESKTOP
}

// Hashes the given string using bcrypt
func HashString(key string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func GetCurrentTime() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now().UTC())
}
