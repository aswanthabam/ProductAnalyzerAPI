package middlewares

import (
	"net/http"
	"productanalyzer/api/db"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
	response "productanalyzer/api/utils/response"
	"productanalyzer/api/websockets"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WebsocketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := websockets.UpgradeConnection(c)
		if err != nil {
			response.SendFailureResponse(c, api_error.NewAPIError("Failed to upgrade connection", http.StatusInternalServerError, "Failed to upgrade connection"))
			c.Abort()
			return
		}
		c.Set("websocket", conn)
		c.Next()
		conn.LoopConnection()
		conn.Close()
	}
}

// WebsocketAuthMiddleware is a middleware to authenticate the user,
// Difference between AuthMiddleware and WebsocketAuthMiddleware is that WebsocketAuthMiddleware is designed for websocket connections.
func WebsocketAuthMiddleware(requireVerifiedEmail bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, exists := c.Get("websocket")
		if !exists {
			response.SendFailureResponse(c, api_error.NewAPIError("Unexpected Error", http.StatusInternalServerError, "Unexpected Error"))
			c.Abort()
			return
		}
		wsConn := conn.(*websockets.WebsocketConnection)
		authError := api_error.NewAPIError("Unauthorized", http.StatusUnauthorized, "Authorization Failed")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			wsConn.SendErrorMesage("Authorization header not found")
			wsConn.Close()
			c.Abort()
			return
		}
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			authError.Message = "Invalid Authorization header format"
			wsConn.SendErrorMesage("Invalid Authorization header format")
			wsConn.Close()
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(bearerToken[1])
		if err != nil {
			wsConn.SendErrorMesage(err.Message)
			wsConn.Close()
			c.Abort()
			return
		}
		userID, err2 := primitive.ObjectIDFromHex(claims.UserID)
		if err2 != nil {
			wsConn.SendErrorMesage("Invalid User ID")
			wsConn.Close()
			c.Abort()
			return
		}
		user, err := db.UserRepository.GetUserByID(userID)
		if err != nil {
			wsConn.SendErrorMesage("User not found")
			wsConn.Close()
			c.Abort()
			return
		}
		if requireVerifiedEmail && !user.EmailVerified {
			wsConn.SendErrorMesage("Email Verification Required")
			wsConn.Close()
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
