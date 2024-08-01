package products_route

type CreateProductResponse struct {
	ProductId string `json:"product_id"`
}

type CreateAccessKeyResponse struct {
	ProductID string `json:"product_id"`
	Scope     string `json:"scope"`
	AccessKey string `json:"access_key"`
}

type ProductInfoResponse struct {
	ID          string `json:"id"` // primary key
	ProductID   string `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BaseUrl     string `json:"base_url"`
}

type ProductAccessKeyResponse struct {
	AccessKey string `json:"access_key"`
	Scope     string `json:"scope"`
	CreatedAt string `json:"created_at"`
}
