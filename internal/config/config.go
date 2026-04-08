package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"test-task/internal/utils"
)

const (
	defaultHTTPAddr       = ":8080"
	defaultBaseURL        = "https://api.green-api.com"
	defaultRequestTimeout = 15 * time.Second
)

type Config struct {
	HTTPAddr       string
	BaseURL        string
	RequestTimeout time.Duration
}

func Load() (Config, error) {
	baseURL := strings.TrimRight(utils.GetEnv("GREEN_API_BASE_URL", defaultBaseURL), "/")
	if err := validateBaseURL(baseURL); err != nil {
		return Config{}, err
	}

	requestTimeout, err := utils.LoadDurationEnv("GREEN_API_REQUEST_TIMEOUT", defaultRequestTimeout)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		HTTPAddr:       utils.GetEnv("HTTP_ADDR", defaultHTTPAddr),
		BaseURL:        baseURL,
		RequestTimeout: requestTimeout,
	}
	return cfg, nil
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
