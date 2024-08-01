package products

import (
	products_db "productanalyzer/api/db/products"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	productsRoute := router.Group("/:product_id")
	{
		protectedProductsRoute := productsRoute.Group("/")
		protectedProductsRoute.Use(middlewares.AccessKeyMiddleware(products_db.PRODUCT_ACCESS_KEY_SCOPE_VISIT))

		{
			protectedProductsRoute.GET("/visit", VisitProduct)
		}

	}
}
