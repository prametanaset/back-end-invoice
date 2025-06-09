package middleware

import (
	"github.com/gofiber/fiber/v2"
	"invoice_project/pkg/apperror"
)

// ErrorResponse formats an error message for JSON responses.
func ErrorResponse(message string) fiber.Map {
	return fiber.Map{"error": message}
}

// ErrorHandler is a Fiber error handler returning JSON using ErrorResponse.
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	switch e := err.(type) {
	case *apperror.StatusError:
		code = e.Code
	case *fiber.Error:
		code = e.Code
	}

	msg := apperror.StatusMessage(code)
	return c.Status(code).JSON(ErrorResponse(msg))
}
