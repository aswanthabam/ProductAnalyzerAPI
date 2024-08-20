package user

import (
	auth_route "productanalyzer/api/api/user/auth"
	profile_route "productanalyzer/api/api/user/profile"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	protectedRoute := router.Group("/")
	protectedRoute.Use(middlewares.AuthMiddleware(false))
	emailProtectedRoute := protectedRoute.Group("/")
	emailProtectedRoute.Use(middlewares.AuthMiddleware(true))

	authRoute := router.Group("/auth")
	protectedAuthRoute := protectedRoute.Group("/auth")
	{
		authRoute.POST("/register", auth_route.Register)
		authRoute.POST("/login", auth_route.Login)
		authRoute.POST("/get-access-token", auth_route.GetAccessToken)

		protectedAuthRoute.POST("/verify-email", auth_route.VerifyEmail)
		protectedAuthRoute.POST("/resend-otp", auth_route.ResendOTP)
	}
	// profileRoute := router.Group("/profile")
	protectedProfileRoute := protectedRoute.Group("/profile")
	{
		protectedProfileRoute.GET("/info", profile_route.GetUserInfo)
	}
}
