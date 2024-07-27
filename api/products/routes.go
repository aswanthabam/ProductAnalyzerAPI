package products

import (
	dashboard_route "productanalyzer/api/api/products/dashboard"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	productsRoute := router.Group("/products")
	{
		protectedProductsRoute := productsRoute.Group("/")
		protectedProductsRoute.Use(middlewares.AuthMiddleware(false))

		protectedProductsRoute.POST("/create", dashboard_route.CreateProduct)
	}
}
