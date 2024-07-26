package middlewares

import (
	"net/http"
	"productanalyzer/api/db"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
		user := db.User{}
		if err := db.Connection.User.FindOne(c, bson.M{"_id": userID}).Decode(&user); err != nil {
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
