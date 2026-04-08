package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"test-task/internal/greenapi"
)

func TestSendMessageValidationError(t *testing.T) {
	app := NewApp(stubFactory{})

	req := httptest.NewRequest("POST", "/api/v1/messages/text", bytes.NewBufferString(`{
		"idInstance":"1101",
		"apiTokenInstance":"token-1",
		"chatId":"79991234567@c.us",
		"message":""
	}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	assertJSONError(t, resp, 400, "validation_error")
}

func TestSendMessageRejectsTrailingJSON(t *testing.T) {
	app := NewApp(stubFactory{})

	req := httptest.NewRequest(
		"POST",
		"/api/v1/messages/text",
		bytes.NewBufferString(`{"idInstance":"1101","apiTokenInstance":"token-1","chatId":"79991234567@c.us","message":"hello"}{"extra":true}`),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	assertJSONError(t, resp, 400, "invalid_request")
}

func TestGetStateSuccess(t *testing.T) {
	app := NewApp(stubFactory{
		api: stubAPI{
			getStateResponse: greenapi.StateResponse{StateInstance: "authorized"},
		},
	})

	req := httptest.NewRequest("POST", "/api/v1/state", bytes.NewBufferString(`{
		"idInstance":"1101",
		"apiTokenInstance":"token-1"
	}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var body greenapi.StateResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.StateInstance != "authorized" {
		t.Fatalf("stateInstance = %q, want %q", body.StateInstance, "authorized")
	}
}

func TestGetStateUpstreamError(t *testing.T) {
	app := NewApp(stubFactory{
		api: stubAPI{
			getStateError: &greenapi.UpstreamError{StatusCode: 500, Body: `{"error":"boom"}`},
		},
	})

	req := httptest.NewRequest("POST", "/api/v1/state", bytes.NewBufferString(`{
		"idInstance":"1101",
		"apiTokenInstance":"token-1"
	}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	assertJSONError(t, resp, 502, "upstream_error")
}

func assertJSONError(t *testing.T, resp *http.Response, wantStatus int, wantCode string) {
	t.Helper()

	defer resp.Body.Close()

	if resp.StatusCode != wantStatus {
		t.Fatalf("status = %d, want %d", resp.StatusCode, wantStatus)
	}

	var payload errorResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if payload.Error != wantCode {
		t.Fatalf("error code = %q, want %q", payload.Error, wantCode)
	}
	if payload.Details == "" {
		t.Fatal("error details must not be empty")
	}
}

type stubFactory struct {
	api      greenapi.API
	withErr  error
	lastCred greenapi.Credentials
}

func (f stubFactory) WithCredentials(creds greenapi.Credentials) (greenapi.API, error) {
	if f.withErr != nil {
		return nil, f.withErr
	}
	return f.api, nil
}

type stubAPI struct {
	getSettingsResponse greenapi.SettingsResponse
	getSettingsError    error
	getStateResponse    greenapi.StateResponse
	getStateError       error
	sendMessageResponse greenapi.SendMessageResponse
	sendMessageError    error
	sendFileResponse    greenapi.SendMessageResponse
	sendFileError       error
}

func (s stubAPI) GetSettings(context.Context) (greenapi.SettingsResponse, error) {
	return s.getSettingsResponse, s.getSettingsError
}

func (s stubAPI) GetStateInstance(context.Context) (greenapi.StateResponse, error) {
	return s.getStateResponse, s.getStateError
}

func (s stubAPI) SendMessage(context.Context, greenapi.SendMessageRequest) (greenapi.SendMessageResponse, error) {
	return s.sendMessageResponse, s.sendMessageError
}

func (s stubAPI) SendFileByURL(context.Context, greenapi.SendFileByURLRequest) (greenapi.SendMessageResponse, error) {
	return s.sendFileResponse, s.sendFileError
}

var _ greenapi.Factory = stubFactory{}
var _ greenapi.API = stubAPI{}

func TestCredentialsValidationError(t *testing.T) {
	app := NewApp(stubFactory{
		withErr: errors.New("idInstance is required"),
	})

	req := httptest.NewRequest("POST", "/api/v1/state", bytes.NewBufferString(`{
		"idInstance":"",
		"apiTokenInstance":"token-1"
	}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	assertJSONError(t, resp, 400, "validation_error")
}
