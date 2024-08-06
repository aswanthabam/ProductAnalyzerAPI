package products

type ProductVisitResponse struct {
	VisitId    string `json:"visit_id"`
	SessionId  string `json:"session_id"`
	LocationId string `json:"location_id"`
}
