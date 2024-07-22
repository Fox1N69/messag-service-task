package dto

// Response represents a generic response structure for API responses.
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SeccessResponse represents a generic response structure for API success reponses without code.
type SuccessResponse struct {
	Message string `json:"message"`
}
