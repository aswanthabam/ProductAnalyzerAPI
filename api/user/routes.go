package user

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", Register)
	}
}
