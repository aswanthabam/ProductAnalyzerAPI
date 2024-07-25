package api

import (
	"productanalyzer/api/api/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	user.SetupRoutes(userRoutes)
}
