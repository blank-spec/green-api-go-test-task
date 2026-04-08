package greenapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"test-task/internal/config"
)

func TestClientGetStateInstance(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/waInstance1101/getStateInstance/token-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(StateResponse{StateInstance: "authorized"})
	}))
	defer ts.Close()

	client := NewClient(config.Config{
		BaseURL:        ts.URL,
		RequestTimeout: time.Second,
	})

	api, err := client.WithCredentials(Credentials{
		IDInstance: "1101",
		APIToken:   "token-1",
	})
	if err != nil {
		t.Fatalf("WithCredentials() error = %v", err)
	}

	resp, err := api.GetStateInstance(context.Background())
	if err != nil {
		t.Fatalf("GetStateInstance() error = %v", err)
	}
	if resp.StateInstance != "authorized" {
		t.Fatalf("unexpected state: %s", resp.StateInstance)
	}
}

func TestClientSendMessageUpstreamError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
	}))
	defer ts.Close()

	client := NewClient(config.Config{
		BaseURL:        ts.URL,
		RequestTimeout: time.Second,
	})

	api, err := client.WithCredentials(Credentials{
		IDInstance: "1101",
		APIToken:   "token-1",
	})
	if err != nil {
		t.Fatalf("WithCredentials() error = %v", err)
	}

	_, err = api.SendMessage(context.Background(), SendMessageRequest{
		ChatID:  "79991234567@c.us",
		Message: "hello",
	})
	if err == nil {
		t.Fatal("expected upstream error")
	}

	upstreamErr, ok := err.(*UpstreamError)
	if !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
	if upstreamErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", upstreamErr.StatusCode)
	}
}

func TestClientRequiresCredentialsBeforeRequest(t *testing.T) {
	client := NewClient(config.Config{
		BaseURL:        "https://example.com",
		RequestTimeout: time.Second,
	})

	_, err := client.GetStateInstance(context.Background())
	if err == nil {
		t.Fatal("expected configuration error")
	}
	if !errors.Is(err, ErrMissingIDInstance) {
		t.Fatalf("unexpected error: %v", err)
	}
}
