package dashboard_route

type CreateProductRequest struct {
	Name        string `form:"name" binding:"required,min=3,max=50"`
	Description string `form:"description" binding:"required,min=3,max=100"`
	BaseUrl     string `form:"base_url" binding:"required,url"`
	ProductID   string `form:"product_id" binding:"required,min=3,max=50"`
}

type CreateAccessKeyRequest struct {
	ProductID string `form:"product_id" binding:"required"`
	Scope     string `form:"scope" binding:"required"`
}

type ProductInfoRequest struct {
	ProductID string `form:"product_id" binding:"required"`
}

type ProductAccessKeysRequest struct {
	ProductID string `form:"product_id" binding:"required"`
}

type DeleteProductRequest struct {
	InstanceId string `form:"instance_id" binding:"required"`
	Type       string `form:"type" binding:"required"`
}
