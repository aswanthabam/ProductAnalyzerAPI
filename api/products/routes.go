package products

import (
	dashboard_route "productanalyzer/api/api/products/dashboard"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	dashboardRoute := router.Group("/dashboard")
	{
		protectedProductsRoute := dashboardRoute.Group("/")
		protectedProductsRoute.Use(middlewares.AuthMiddleware(false))

		protectedProductsRoute.POST("/create", dashboard_route.CreateProduct)
	}
}
