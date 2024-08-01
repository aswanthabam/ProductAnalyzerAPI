package dashboard_route

import (
	products_db "productanalyzer/api/db/products"
	user_db "productanalyzer/api/db/user"
	api_error "productanalyzer/api/errors"
	"productanalyzer/api/utils"
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
	productId, err := products_db.CreateProduct(&product)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "Product created successfully", CreateProductResponse{
		ProductId: productId.Hex(),
		AccessKey: product.AccessKeys[0].AccessKey,
	}, nil)
}

// Create Access Key Request [POST]
func CreateAccessKey(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params CreateAccessKeyRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := products_db.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	key, err2 := utils.GenerateAPIKey()
	if err2 != nil {
		response.SendFailureResponse(c, api_error.UnexpectedError(err))
		return
	}
	if params.Scope != products_db.PRODUCT_ACCESS_KEY_SCOPE_VISIT && params.Scope != products_db.PRODUCT_ACCESS_KEY_SCOPE_ALL {
		response.SendFailureResponse(c, api_error.NewAPIError("Invalid Scope", 400, "The given scope is invalid"))
		return
	}
	accessKey := products_db.ProductAccessKey{
		AccessKey: key,
		Scope:     params.Scope,
		CreatedAt: utils.GetCurrentTime(),
	}
	product.AccessKeys = append(product.AccessKeys, accessKey)
	err = products_db.UpdateProduct(*product)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "Access Key created successfully", CreateAccessKeyResponse{
		AccessKey: key,
		ProductID: product.ProductID,
		Scope:     params.Scope,
	}, nil)
}

// Product Info Request [POST]
func ProductInfo(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params ProductInfoRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := products_db.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	response.SendSuccessResponse(c, "Product information", ProductInfoResponse{
		Name:        product.Name,
		Description: product.Description,
		BaseUrl:     product.BaseUrl,
		ProductID:   product.ProductID,
	}, nil)
}

func ProductAccessKeys(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		response.SendFailureResponse(c, api_error.UnexpectedError(nil))
		return
	}
	user := usr.(*user_db.User)
	var params ProductAccessKeysRequest
	if err := c.ShouldBind(&params); err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	product, err := products_db.GetProductByProductIDAUserID(params.ProductID, user.ID)
	if err != nil {
		response.SendFailureResponse(c, err)
		return
	}
	keys := []ProductAccessKeyResponse{}
	for _, key := range product.AccessKeys {
		keys = append(keys, ProductAccessKeyResponse{
			AccessKey: key.AccessKey,
			Scope:     key.Scope,
			CreatedAt: key.CreatedAt.Time().UTC().String(),
		})
	}
	response.SendSuccessResponse(c, "Product access keys", keys, nil)
}
