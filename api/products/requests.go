package products

type VisitProductRequest struct {
	From   string `form:"from"`
	Page   string `form:"page"`
	Method string `form:"method"`
}
