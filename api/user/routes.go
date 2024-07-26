package user

import (
	"productanalyzer/api/api/user/auth"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	authRoute := router.Group("/auth")
	{
		protectedAuthRoute := authRoute.Group("/")
		protectedAuthRoute.Use(middlewares.AuthMiddleware(false))

		authRoute.POST("/register", auth.Register)
		protectedAuthRoute.POST("/verify-email", auth.VerifyEmail)
	}
}
