package models

// Response represents a generic response structure for API responses.
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
