package products_route

import "time"

const (
	VISIT_LOG_ENTRY_UNIT_MINUTE = "minute"
	VISIT_LOG_ENTRY_UNIT_HOUR   = "hour"
	VISIT_LOG_ENTRY_UNIT_DAY    = "day"
	VISIT_LOG_ENTRY_UNIT_MONTH  = "month"
)

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

type VisitLogRequest struct {
	ProductID string    `form:"product_id" binding:"required"`
	FromDate  time.Time `form:"from_date" binding:"required" time_format:"2006-01-02 15:04:05"`
	ToDate    time.Time `form:"to_date" time_format:"2006-01-02 15:04:05" default:"now"`
	Unit      string    `form:"unit" binding:"required"`
}
