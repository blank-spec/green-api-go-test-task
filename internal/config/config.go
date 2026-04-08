package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	defaultHTTPAddr       = ":8080"
	defaultBaseURL        = "https://api.green-api.com"
	defaultRequestTimeout = 15 * time.Second
)

// Config contains runtime settings for the HTTP server and Green API client.
type Config struct {
	HTTPAddr       string
	BaseURL        string
	RequestTimeout time.Duration
}

// Load reads configuration from environment variables and validates it.
func Load() (Config, error) {
	baseURL := strings.TrimRight(getEnv("GREEN_API_BASE_URL", defaultBaseURL), "/")
	if err := validateBaseURL(baseURL); err != nil {
		return Config{}, err
	}

	requestTimeout, err := loadDurationEnv("GREEN_API_REQUEST_TIMEOUT", defaultRequestTimeout)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		HTTPAddr:       getEnv("HTTP_ADDR", defaultHTTPAddr),
		BaseURL:        baseURL,
		RequestTimeout: requestTimeout,
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func loadDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return duration, nil
}

func validateBaseURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("parse GREEN_API_BASE_URL: %w", err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("GREEN_API_BASE_URL must use http or https")
	}
	if parsedURL.Host == "" {
		return fmt.Errorf("GREEN_API_BASE_URL must include a host")
	}

	return nil
}
