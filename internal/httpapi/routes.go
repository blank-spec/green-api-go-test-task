package httpapi

import "github.com/gofiber/fiber/v2"

func registerRoutes(app *fiber.App, h *handler) {
	app.Get("/", h.index)
	app.Get("/healthz", h.health)
	app.Post("/api/v1/settings", h.postSettings)
	app.Post("/api/v1/state", h.postState)
	app.Post("/api/v1/messages/text", h.sendMessage)
	app.Post("/api/v1/messages/file", h.sendFileByURL)
}
