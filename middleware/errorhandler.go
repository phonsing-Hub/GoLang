package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/utils/response"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := err.Error()

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}

	return c.Status(code).JSON(response.ErrorResponse{
		Success: false,
		Data:    nil,
		Error: &response.ErrorDetail{
			Code:    http.StatusText(code),
			Message: msg,
		},
	})

}
