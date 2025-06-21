package response

import "encoding/json"

type SWListDataDetail struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Data  json.RawMessage `json:"data" swaggertype:"array,object"`
}

type SWListSuccessResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    *SWListDataDetail `json:"data"`
	Error   json.RawMessage   `json:"error" swaggertype:"object"`
}

type SWSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    json.RawMessage `json:"data" swaggertype:"object"`
	Error   json.RawMessage `json:"error" swaggertype:"object"`
}

type SWErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SWErrorResponse struct {
	Success bool           `json:"success" example:"false"`
	Data    any            `json:"data"`
	Error   *SWErrorDetail `json:"error"`
}
