package httpapi

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"test-task/internal/greenapi"
)

type requestError struct {
	status  int
	code    string
	details string
}

func (e *requestError) Error() string {
	return e.details
}

func newRequestError(status int, code string, err error) error {
	return &requestError{
		status:  status,
		code:    code,
		details: err.Error(),
	}
}

func wrapClientError(err error) error {
	var upstreamErr *greenapi.UpstreamError
	if errors.As(err, &upstreamErr) {
		details := strings.TrimSpace(upstreamErr.Body)
		if details == "" {
			details = upstreamErr.Error()
		}

		return &requestError{
			status:  fiber.StatusBadGateway,
			code:    "upstream_error",
			details: details,
		}
	}

	return &requestError{
		status:  fiber.StatusInternalServerError,
		code:    "internal_error",
		details: err.Error(),
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	var reqErr *requestError
	if errors.As(err, &reqErr) {
		return c.Status(reqErr.status).JSON(errorResponse{
			Error:   reqErr.code,
			Details: reqErr.details,
		})
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(errorResponse{
			Error:   "http_error",
			Details: fiberErr.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{
		Error:   "internal_error",
		Details: err.Error(),
	})
}
