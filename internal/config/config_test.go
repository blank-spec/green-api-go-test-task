package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("GREEN_API_BASE_URL", "https://example.com/")
	t.Setenv("GREEN_API_REQUEST_TIMEOUT", "7s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":9090")
	}
	if cfg.BaseURL != "https://example.com" {
		t.Fatalf("BaseURL = %q, want %q", cfg.BaseURL, "https://example.com")
	}
	if cfg.RequestTimeout != 7*time.Second {
		t.Fatalf("RequestTimeout = %v, want %v", cfg.RequestTimeout, 7*time.Second)
	}
}

func TestLoadInvalidTimeout(t *testing.T) {
	t.Setenv("GREEN_API_REQUEST_TIMEOUT", "abc")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}

func TestLoadInvalidBaseURL(t *testing.T) {
	t.Setenv("GREEN_API_BASE_URL", "://bad-url")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}
