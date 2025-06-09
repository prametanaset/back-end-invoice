package apperror

import "github.com/gofiber/fiber/v2"

// statusText holds application specific text for common HTTP status codes.
var statusText = map[int]string{
	fiber.StatusBadRequest:          "BAD_REQUEST",
	fiber.StatusUnauthorized:        "UNAUTHORIZED",
	fiber.StatusNotFound:            "NOT_FOUND",
	fiber.StatusInternalServerError: "INTERNAL_SERVER_ERROR",
}

// StatusMessage returns the application error message for the given status code.
func StatusMessage(code int) string {
	if msg, ok := statusText[code]; ok {
		return msg
	}
	return "UNKNOWN_ERROR"
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct{ Code int }

// Error implements the error interface using the mapped status message.
func (e *StatusError) Error() string { return StatusMessage(e.Code) }

// New creates a new StatusError with the provided status code.
func New(code int) error { return &StatusError{Code: code} }
