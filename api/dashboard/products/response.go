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

type VisitLogResponse struct {
	Logs            []VisitLogEntry `json:"logs"`
	TotalSessions   int64           `json:"total_sessions"`
	TotalActivities int64           `json:"total_activities"`
}

type VisitLogEntry struct {
	Type          string `json:"type"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ActivityCount int64  `json:"count"`
	SessionCount  int64  `json:"session_count"`
	Referer       string `json:"referer"`
}
