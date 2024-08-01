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
		protectedProductsRoute.Use(middlewares.AuthMiddleware(true))

		protectedProductsRoute.POST("/create", dashboard_route.CreateProduct)
		protectedProductsRoute.POST("/create-access-key", dashboard_route.CreateAccessKey)
		protectedProductsRoute.POST("/info", dashboard_route.ProductInfo)
		protectedProductsRoute.POST("/access-keys", dashboard_route.ProductAccessKeys)
	}
}
