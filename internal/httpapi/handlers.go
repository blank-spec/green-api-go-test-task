package httpapi

import (
	"github.com/gofiber/fiber/v2"

	"test-task/internal/greenapi"
)

type handler struct {
	factory greenapi.Factory
}

func newHandler(factory greenapi.Factory) *handler {
	return &handler{factory: factory}
}

func (h *handler) health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(healthResponse{Status: "ok"})
}

func (h *handler) postSettings(c *fiber.Ctx) error {
	var req settingsRequest
	if err := decodeJSON(c, &req); err != nil {
		return newRequestError(fiber.StatusBadRequest, "invalid_request", err)
	}

	client, err := h.withCredentials(req.credentialsRequest)
	if err != nil {
		return newRequestError(fiber.StatusBadRequest, "validation_error", err)
	}

	resp, err := client.GetSettings(c.UserContext())
	if err != nil {
		return wrapClientError(err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) postState(c *fiber.Ctx) error {
	var req stateRequest
	if err := decodeJSON(c, &req); err != nil {
		return newRequestError(fiber.StatusBadRequest, "invalid_request", err)
	}

	client, err := h.withCredentials(req.credentialsRequest)
	if err != nil {
		return newRequestError(fiber.StatusBadRequest, "validation_error", err)
	}

	resp, err := client.GetStateInstance(c.UserContext())
	if err != nil {
		return wrapClientError(err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) sendMessage(c *fiber.Ctx) error {
	var req sendMessagePageRequest
	if err := decodeJSON(c, &req); err != nil {
		return newRequestError(fiber.StatusBadRequest, "invalid_request", err)
	}

	client, err := h.withCredentials(req.credentialsRequest)
	if err != nil {
		return newRequestError(fiber.StatusBadRequest, "validation_error", err)
	}
	if err := ValidateSendMessage(req.SendMessageRequest); err != nil {
		return newRequestError(fiber.StatusBadRequest, "validation_error", err)
	}

	resp, err := client.SendMessage(c.UserContext(), req.SendMessageRequest)
	if err != nil {
		return wrapClientError(err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) sendFileByURL(c *fiber.Ctx) error {
	var req sendFilePageRequest
	if err := decodeJSON(c, &req); err != nil {
		return newRequestError(fiber.StatusBadRequest, "invalid_request", err)
	}

	client, err := h.withCredentials(req.credentialsRequest)
	if err != nil {
		return newRequestError(fiber.StatusBadRequest, "validation_error", err)
	}
	if err := ValidateSendFileByURL(req.SendFileByURLRequest); err != nil {
		return newRequestError(fiber.StatusBadRequest, "validation_error", err)
	}

	resp, err := client.SendFileByURL(c.UserContext(), req.SendFileByURLRequest)
	if err != nil {
		return wrapClientError(err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) withCredentials(req credentialsRequest) (greenapi.API, error) {
	return h.factory.WithCredentials(greenapi.Credentials{
		IDInstance: req.IDInstance,
		APIToken:   req.APIToken,
	})
}
