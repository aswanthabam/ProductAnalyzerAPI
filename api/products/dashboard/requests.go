package dashboard_route

type CreateProductRequest struct {
	Name        string `form:"name" binding:"required,min=3,max=50"`
	Description string `form:"description" binding:"required,min=3,max=100"`
	BaseUrl     string `form:"base_url" binding:"required,url"`
	ProductID   string `form:"product_id" binding:"required,min=3,max=50"`
}
