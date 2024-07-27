package api

import (
	"productanalyzer/api/api/products"
	"productanalyzer/api/api/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	productsRoute := router.Group("/products")

	user.SetupRoutes(userRoutes)
	products.SetupRoutes(productsRoute)
}
