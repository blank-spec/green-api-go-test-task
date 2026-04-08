package httpapi

import (
	"embed"

	"github.com/gofiber/fiber/v2"
)

//go:embed assets/index.html
var assets embed.FS

func (h *handler) index(c *fiber.Ctx) error {
	page, err := assets.ReadFile("assets/index.html")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load page")
	}

	c.Type("html", "utf-8")
	return c.Send(page)
}
