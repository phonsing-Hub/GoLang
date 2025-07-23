package response

type SWSuccessResponse struct {
	Success bool     `json:"success" example:"true"`
	Data    struct{} `json:"data"`
	Error   struct{} `json:"error"`
}

type SWListSuccessResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    *SWListDataDetail `json:"data"`
	Error   struct{}          `json:"error"`
}

type SWListDataDetail struct {
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Data  []struct{} `json:"data"`
}

type SWErrorResponse struct {
	Success bool           `json:"success" example:"false"`
	Data    struct{}       `json:"data"`
	Error   *SWErrorDetail `json:"error"`
}

type SWErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
