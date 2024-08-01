package products

import (
	products_db "productanalyzer/api/db/products"
	api_error "productanalyzer/api/errors"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
)

func VisitProduct(c *gin.Context) {
	product, exists := c.Get("product")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	product = product.(*products_db.Product)
	var params VisitProductRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "Product visited successfully", nil, nil)
}
