package utils

import (
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	t.Setenv("TEST_ENV_VALUE", " value ")

	value := GetEnv("TEST_ENV_VALUE", "fallback")
	if value != "value" {
		t.Fatalf("GetEnv() = %q, want %q", value, "value")
	}
}

func TestGetEnvFallback(t *testing.T) {
	value := GetEnv("MISSING_ENV_VALUE", "fallback")
	if value != "fallback" {
		t.Fatalf("GetEnv() = %q, want %q", value, "fallback")
	}
}

func TestLoadDurationEnv(t *testing.T) {
	t.Setenv("TEST_DURATION", "7s")

	value, err := LoadDurationEnv("TEST_DURATION", 3*time.Second)
	if err != nil {
		t.Fatalf("LoadDurationEnv() error = %v", err)
	}
	if value != 7*time.Second {
		t.Fatalf("LoadDurationEnv() = %v, want %v", value, 7*time.Second)
	}
}

func TestLoadDurationEnvFallback(t *testing.T) {
	value, err := LoadDurationEnv("MISSING_DURATION", 3*time.Second)
	if err != nil {
		t.Fatalf("LoadDurationEnv() error = %v", err)
	}
	if value != 3*time.Second {
		t.Fatalf("LoadDurationEnv() = %v, want %v", value, 3*time.Second)
	}
}

func TestLoadDurationEnvInvalid(t *testing.T) {
	t.Setenv("TEST_DURATION", "bad")

	if _, err := LoadDurationEnv("TEST_DURATION", 3*time.Second); err == nil {
		t.Fatal("LoadDurationEnv() error = nil, want error")
	}
}
