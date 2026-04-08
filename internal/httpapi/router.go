package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"test-task/internal/greenapi"
)

const maxRequestBodyBytes = 1 << 20

type handler struct {
	factory greenapi.Factory
}

type credentialsRequest struct {
	IDInstance string `json:"idInstance"`
	APIToken   string `json:"apiTokenInstance"`
}

type settingsRequest struct {
	credentialsRequest
}

type stateRequest struct {
	credentialsRequest
}

type sendMessagePageRequest struct {
	credentialsRequest
	greenapi.SendMessageRequest
}

type sendFilePageRequest struct {
	credentialsRequest
	greenapi.SendFileByURLRequest
}

type errorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func NewRouter(factory greenapi.Factory) *fiber.App {
	h := &handler{factory: factory}

	app := fiber.New(fiber.Config{
		AppName:   "green-api-form",
		BodyLimit: maxRequestBodyBytes,
	})
	app.Use(loggingMiddleware)

	app.Get("/", h.index)
	app.Get("/healthz", h.health)
	app.Post("/api/v1/settings", h.postSettings)
	app.Post("/api/v1/state", h.postState)
	app.Post("/api/v1/messages/text", h.sendMessage)
	app.Post("/api/v1/messages/file", h.sendFileByURL)

	return app
}

func (h *handler) health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(map[string]string{"status": "ok"})
}

func (h *handler) postSettings(c *fiber.Ctx) error {
	var req settingsRequest
	if err := decodeJSON(c, &req); err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "invalid_request", err)
	}

	client, err := h.withCredentials(req.credentialsRequest)
	if err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "validation_error", err)
	}

	resp, err := client.GetSettings(c.UserContext())
	if err != nil {
		return writeClientError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) postState(c *fiber.Ctx) error {
	var req stateRequest
	if err := decodeJSON(c, &req); err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "invalid_request", err)
	}

	client, err := h.withCredentials(req.credentialsRequest)
	if err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "validation_error", err)
	}

	resp, err := client.GetStateInstance(c.UserContext())
	if err != nil {
		return writeClientError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) sendMessage(c *fiber.Ctx) error {
	var req sendMessagePageRequest
	if err := decodeJSON(c, &req); err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "invalid_request", err)
	}

	client, err := h.withCredentials(req.credentialsRequest)
	if err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "validation_error", err)
	}
	if err := ValidateSendMessage(req.SendMessageRequest); err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "validation_error", err)
	}

	resp, err := client.SendMessage(c.UserContext(), req.SendMessageRequest)
	if err != nil {
		return writeClientError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *handler) sendFileByURL(c *fiber.Ctx) error {
	var req sendFilePageRequest
	if err := decodeJSON(c, &req); err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "invalid_request", err)
	}

	client, err := h.withCredentials(req.credentialsRequest)
	if err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "validation_error", err)
	}
	if err := ValidateSendFileByURL(req.SendFileByURLRequest); err != nil {
		return writeJSONError(c, fiber.StatusBadRequest, "validation_error", err)
	}

	resp, err := client.SendFileByURL(c.UserContext(), req.SendFileByURLRequest)
	if err != nil {
		return writeClientError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func decodeJSON(c *fiber.Ctx, out any) error {
	body := c.Body()
	if len(body) == 0 {
		return errors.New("request body is required")
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != nil && !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

func (h *handler) withCredentials(req credentialsRequest) (greenapi.API, error) {
	return h.factory.WithCredentials(greenapi.Credentials{
		IDInstance: req.IDInstance,
		APIToken:   req.APIToken,
	})
}

func writeClientError(c *fiber.Ctx, err error) error {
	var upstreamErr *greenapi.UpstreamError
	if errors.As(err, &upstreamErr) {
		details := strings.TrimSpace(upstreamErr.Body)
		if details == "" {
			details = upstreamErr.Error()
		}
		return writeJSONError(c, fiber.StatusBadGateway, "upstream_error", errors.New(details))
	}

	return writeJSONError(c, fiber.StatusInternalServerError, "internal_error", err)
}

func writeJSONError(c *fiber.Ctx, status int, code string, err error) error {
	return c.Status(status).JSON(errorResponse{
		Error:   code,
		Details: err.Error(),
	})
}

func loggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	log.Printf("%s %s %d %s", c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start))
	return err
}
