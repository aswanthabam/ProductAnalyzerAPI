package dashboard_route

import (
	products_db "productanalyzer/api/db/products"
	user_db "productanalyzer/api/db/user"
	api_error "productanalyzer/api/errors"
	response "productanalyzer/api/utils/response"

	"github.com/gin-gonic/gin"
)

// Create Product Request [POST]
func CreateProduct(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params CreateProductRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product := products_db.Product{
		Name:        params.Name,
		Description: params.Description,
		BaseUrl:     params.BaseUrl,
		ProductID:   params.ProductID,
		AccessKeys:  []products_db.ProductAccessKey{},
		UserID:      user.ID,
	}
	productId, err := products_db.CreateProduct(product)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "Product created successfully", CreateProductResponse{
		ProductId: productId.Hex(),
	}, nil)
}
