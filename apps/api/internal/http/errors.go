package http

import "github.com/gofiber/fiber/v2"

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func writeError(c *fiber.Ctx, status int, code, message string, fields map[string]string) error {
	return c.Status(status).JSON(errorEnvelope{
		Error: apiError{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
}
