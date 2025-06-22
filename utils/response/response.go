package response

import (
	"github.com/gofiber/fiber/v2"
)

type SuccessResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
	Error   any  `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool         `json:"success"`
	Data    any          `json:"data"`
	Error   *ErrorDetail `json:"error"`
}

func OK[T any](c *fiber.Ctx, data T, optionalStatus ...int) error {
	status := fiber.StatusOK

	if len(optionalStatus) > 0 {
		status = optionalStatus[0]
	}

	return c.Status(status).JSON(SuccessResponse[T]{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

func Fail(c *fiber.Ctx, code, message string, status int) error {
	return c.Status(status).JSON(ErrorResponse{
		Success: false,
		Data:    nil,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
