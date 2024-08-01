package api

import (
	"productanalyzer/api/api/dashboard"
	"productanalyzer/api/api/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	dashboardRoutes := router.Group("/dashboard")

	user.SetupRoutes(userRoutes)
	dashboard.SetupRoutes(dashboardRoutes)
}
