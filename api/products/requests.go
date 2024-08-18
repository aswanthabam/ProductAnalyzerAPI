package products

type VisitProductRequest struct {
	From      string `form:"from" json:"from"`
	Page      string `form:"page" json:"page"`
	Method    string `form:"method" json:"method"`
	SessionId string `form:"session_id" json:"session_id"`
}
