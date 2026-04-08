package httpapi

import (
	"github.com/gofiber/fiber/v2"

	"test-task/internal/greenapi"
)

func NewRouter(factory greenapi.Factory) *fiber.App {
	return NewApp(factory)
}
