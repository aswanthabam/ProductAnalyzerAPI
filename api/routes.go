package api

import (
	"productanalyzer/api/api/dashboard"
	"productanalyzer/api/api/products"
	"productanalyzer/api/api/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	dashboardRoutes := router.Group("/dashboard")
	productsRoutes := router.Group("/products")

	user.SetupRoutes(userRoutes)
	dashboard.SetupRoutes(dashboardRoutes)
	products.SetupRoutes(productsRoutes)
}
