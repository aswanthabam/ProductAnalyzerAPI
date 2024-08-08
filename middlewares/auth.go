package middlewares

import (
	"net/http"
	products_db "productanalyzer/api/db/products"
	user_db "productanalyzer/api/db/user"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthMiddleware is a middleware to authenticate the user, it checks the Authorization header and validates the token.
// If requireVerifiedEmail is true, it will also check if the user has verified their email.
// If the user is authenticated, it will set the user object in the context.
func AuthMiddleware(requireVerifiedEmail bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authError := api_error.NewAPIError("Unauthorized", http.StatusUnauthorized, "Authorization Failed")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.SendFailureResponse(c, authError)
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			authError.Message = "Invalid Authorization header format"
			response.SendFailureResponse(c, authError)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(bearerToken[1])
		if err != nil {
			response.SendFailureResponse(c, err)
			c.Abort()
			return
		}
		userID, err2 := primitive.ObjectIDFromHex(claims.UserID)
		if err2 != nil {
			response.SendFailureResponse(c, api_error.NewAPIError("Invalid User ID", http.StatusBadRequest, "Invalid User ID"))
			c.Abort()
			return
		}
		user, err := user_db.GetUserByID(userID)
		if err != nil {
			response.SendFailureResponse(c, api_error.NewAPIError("User not found", http.StatusNotFound, "User not found"))
			c.Abort()
			return
		}
		if requireVerifiedEmail && !user.EmailVerified {
			response.SendFailureResponse(c, api_error.NewAPIError("Email Verification Required", http.StatusExpectationFailed, "email is not verified, check your mail for otp"))
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func AccessKeyMiddleware(scope string, strictHeader bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authError := api_error.NewAPIError("Unauthorized", http.StatusUnauthorized, "Invalid API Key. Please provide a valid API Key in the X-API-Key header")
		var key string
		key = c.GetHeader("X-API-Key")
		if !strictHeader && key == "" {
			key = c.Query("api_key")
		}
		if key == "" {
			response.SendFailureResponse(c, authError)
			c.Abort()
			return
		}
		accessKey, err := products_db.ValidateAPIKey(key, scope)
		if err != nil {
			response.SendFailureResponse(c, err)
			c.Abort()
			return
		}
		c.Set("access_key", accessKey)
		product, err := products_db.GetProductByID(accessKey.ProductID)
		if err != nil {
			response.SendFailureResponse(c, err)
			c.Abort()
			return
		}
		c.Set("product", product)
		c.Next()
	}
}
