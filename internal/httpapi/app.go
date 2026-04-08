package httpapi

import (
	"github.com/gofiber/fiber/v2"

	"test-task/internal/greenapi"
)

const maxRequestBodyBytes = 1 << 20

func NewApp(factory greenapi.Factory) *fiber.App {
	h := newHandler(factory)

	app := fiber.New(fiber.Config{
		AppName:      "green-api-form",
		BodyLimit:    maxRequestBodyBytes,
		ErrorHandler: errorHandler,
	})

	app.Use(loggingMiddleware)
	registerRoutes(app, h)

	return app
}
