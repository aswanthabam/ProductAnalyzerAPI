package products_route

import "time"

const (
	VISIT_LOG_ENTRY_UNIT_MINUTE = "minute"
	VISIT_LOG_ENTRY_UNIT_HOUR   = "hour"
	VISIT_LOG_ENTRY_UNIT_DAY    = "day"
	VISIT_LOG_ENTRY_UNIT_MONTH  = "month"
)

type CreateProductRequest struct {
	Name        string `form:"name" binding:"required,min=3,max=50" json:"name"`
	Description string `form:"description" binding:"required,min=3,max=100" json:"description"`
	BaseUrl     string `form:"base_url" binding:"required,url" json:"base_url"`
	ProductID   string `form:"product_id" binding:"required,min=3,max=50" json:"product_id"`
}

type CreateAccessKeyRequest struct {
	ProductID string `form:"product_id" binding:"required" json:"product_id"`
	Scope     string `form:"scope" binding:"required" json:"scope"`
}

type ProductInfoRequest struct {
	ProductID string `form:"product_id" binding:"required" json:"product_id"`
}

type ProductAccessKeysRequest struct {
	ProductID string `form:"product_id" binding:"required" json:"product_id"`
}

type DeleteProductRequest struct {
	InstanceId string `form:"instance_id" binding:"required" json:"instance_id"`
	Type       string `form:"type" binding:"required" json:"type"`
}

type VisitLogRequest struct {
	ProductID string    `form:"product_id" binding:"required" json:"product_id"`
	FromDate  time.Time `form:"from_date" binding:"required" time_format:"2006-01-02 15:04:05" json:"from_date"`
	ToDate    time.Time `form:"to_date" time_format:"2006-01-02 15:04:05" default:"now" json:"to_date"`
	Unit      string    `form:"unit" binding:"required" json:"unit"`
}
