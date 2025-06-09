package middleware

import "github.com/gofiber/fiber/v2"

// ErrorResponse formats an error message for JSON responses.
func ErrorResponse(message string) fiber.Map {
	return fiber.Map{"error": message}
}

// ErrorHandler is a Fiber error handler returning JSON using ErrorResponse.
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		if e.Message != "" {
			err = e
		}
	}
	return c.Status(code).JSON(ErrorResponse(err.Error()))
}
