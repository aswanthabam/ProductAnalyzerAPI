package products

import (
	db_types "productanalyzer/api/db/types"
	"productanalyzer/api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	protectedProductsRoute := router.Group("/")
	protectedProductsRoute.Use(middlewares.AccessKeyMiddleware(db_types.PRODUCT_ACCESS_KEY_SCOPE_VISIT, true))
	{
		protectedProductsRoute.GET("/log", VisitProduct)
	}
}
