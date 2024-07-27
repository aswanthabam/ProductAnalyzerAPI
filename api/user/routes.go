package user

import (
	auth_route "productanalyzer/api/api/user/auth"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	authRoute := router.Group("/auth")
	{
		protectedAuthRoute := authRoute.Group("/")
		protectedAuthRoute.Use(middlewares.AuthMiddleware(false))

		authRoute.POST("/register", auth_route.Register)
		authRoute.POST("/login", auth_route.Login)
		protectedAuthRoute.POST("/verify-email", auth_route.VerifyEmail)
		protectedAuthRoute.POST("/resend-otp", auth_route.ResendOTP)
	}
}
