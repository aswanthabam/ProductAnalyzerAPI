package user

import (
	"productanalyzer/api/api/user/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	authRoute := router.Group("/auth")
	{
		authRoute.POST("/register", auth.Register)
	}
}
