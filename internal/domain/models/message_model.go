package models

type CreateMessageReq struct {
	Content  string `json:"content" validate:"required"`
	StatusID int64  `json:"status_id" validate:"required"`
}
