package model

type ResponseMeta struct {
	Limit  uint `json:"limit"`
	Offset uint `json:"offset"`
	Total  uint `json:"total"`
}

type DataResponse struct {
	Message string        `json:"message"`
	Data    any           `json:"data,omitempty"`
	Meta    *ResponseMeta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
