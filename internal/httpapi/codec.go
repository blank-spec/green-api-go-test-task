package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
)

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

	if err := consumeTrailingJSON(decoder); err != nil {
		return err
	}

	return nil
}

func consumeTrailingJSON(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err == nil {
		return errors.New("request body must contain a single JSON object")
	} else if !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}
