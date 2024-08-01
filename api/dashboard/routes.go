package dashboard

import (
	products_route "productanalyzer/api/api/dashboard/products"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	productsRoute := router.Group("/products")
	{
		protectedProductsRoute := productsRoute.Group("/")
		protectedProductsRoute.Use(middlewares.AuthMiddleware(true))

		protectedProductsRoute.POST("/create", products_route.CreateProduct)
		protectedProductsRoute.POST("/create-access-key", products_route.CreateAccessKey)
		protectedProductsRoute.GET("/info", products_route.ProductInfo)
		protectedProductsRoute.GET("/access-keys", products_route.ProductAccessKeys)
		protectedProductsRoute.DELETE("/delete", products_route.DeleteProduct)
	}
}
