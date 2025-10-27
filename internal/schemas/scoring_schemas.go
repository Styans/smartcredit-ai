package schemas

type ScoringRequest struct {
	Query string `json:"query" binding:"required,min=5"`
}

type ScoringResponse struct {
	Answer string `json:"answer"`
}