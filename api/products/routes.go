package products

import (
	products_db "productanalyzer/api/db/products"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	protectedProductsRoute := router.Group("/")
	protectedProductsRoute.Use(middlewares.AccessKeyMiddleware(products_db.PRODUCT_ACCESS_KEY_SCOPE_VISIT, true))
	{
		protectedProductsRoute.GET("/log", VisitProduct)
	}
}
