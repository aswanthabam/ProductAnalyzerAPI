package middlewares

import (
	"productanalyzer/api/utils"

	"github.com/gin-gonic/gin"
)

func BotDetectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isBot := utils.IsBot(c.Request)
		c.Set("isBot", isBot)
		// if isBot {
		// 	c.Abort()
		// 	response.SendFailureResponse(c, api_error.NewAPIError("Unauthorized", http.StatusUnauthorized, "Bot Detected"))
		// 	return
		// }
		c.Next()
	}
}
