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
		{
			protectedProductsRoute.POST("/create", products_route.CreateProduct)
			protectedProductsRoute.GET("/list", products_route.ListProducts)
			protectedProductsRoute.POST("/create-access-key", products_route.CreateAccessKey)
			protectedProductsRoute.GET("/info", products_route.ProductInfo)
			protectedProductsRoute.GET("/access-keys", products_route.ProductAccessKeys)
			protectedProductsRoute.DELETE("/delete", products_route.DeleteProduct)

			protectedProductsRoute.GET("/logs", products_route.VisitLog)
		}

		socketRoute := productsRoute.Group("/")
		socketRoute.Use(middlewares.WebsocketMiddleware())
		protectedSocketRoute := socketRoute.Group("/")
		protectedSocketRoute.Use(middlewares.WebsocketAuthMiddleware(true))
		{
			protectedSocketRoute.GET("/logs/ws/:product_id", products_route.HandleLogWebsocket)
		}
	}
}
