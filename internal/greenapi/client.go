package greenapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"test-task/internal/config"
)

var (
	ErrMissingBaseURL    = errors.New("baseURL is required")
	ErrMissingHTTPClient = errors.New("http client is required")
	ErrMissingIDInstance = errors.New("idInstance is required")
	ErrMissingAPIToken   = errors.New("apiTokenInstance is required")
)

type API interface {
	GetSettings(ctx context.Context) (SettingsResponse, error)
	GetStateInstance(ctx context.Context) (StateResponse, error)
	SendMessage(ctx context.Context, reqBody SendMessageRequest) (SendMessageResponse, error)
	SendFileByURL(ctx context.Context, reqBody SendFileByURLRequest) (SendMessageResponse, error)
}

type Factory interface {
	WithCredentials(creds Credentials) (API, error)
}

type Client struct {
	baseURL    string
	idInstance string
	apiToken   string
	httpClient *http.Client
}

type UpstreamError struct {
	StatusCode int
	Body       string
}

func (e *UpstreamError) Error() string {
	return fmt.Sprintf("green api returned status %d", e.StatusCode)
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		baseURL:    cfg.BaseURL,
		httpClient: &http.Client{Timeout: cfg.RequestTimeout},
	}
}

func (c *Client) WithCredentials(creds Credentials) (API, error) {
	if strings.TrimSpace(creds.IDInstance) == "" {
		return nil, ErrMissingIDInstance
	}
	if strings.TrimSpace(creds.APIToken) == "" {
		return nil, ErrMissingAPIToken
	}

	return &Client{
		baseURL:    c.baseURL,
		idInstance: strings.TrimSpace(creds.IDInstance),
		apiToken:   strings.TrimSpace(creds.APIToken),
		httpClient: c.httpClient,
	}, nil
}

func (c *Client) GetSettings(ctx context.Context) (SettingsResponse, error) {
	var resp SettingsResponse
	err := c.do(ctx, http.MethodGet, "getSettings", nil, &resp)
	return resp, err
}

func (c *Client) GetStateInstance(ctx context.Context) (StateResponse, error) {
	var resp StateResponse
	err := c.do(ctx, http.MethodGet, "getStateInstance", nil, &resp)
	return resp, err
}

func (c *Client) SendMessage(ctx context.Context, reqBody SendMessageRequest) (SendMessageResponse, error) {
	var resp SendMessageResponse
	err := c.do(ctx, http.MethodPost, "sendMessage", reqBody, &resp)
	return resp, err
}

func (c *Client) SendFileByURL(ctx context.Context, reqBody SendFileByURLRequest) (SendMessageResponse, error) {
	var resp SendMessageResponse
	err := c.do(ctx, http.MethodPost, "sendFileByUrl", reqBody, &resp)
	return resp, err
}

func (c *Client) do(ctx context.Context, method, endpoint string, payload any, out any) error {
	if err := c.validateConfigured(); err != nil {
		return err
	}

	var body io.Reader
	if payload != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		body = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpointURL(endpoint), body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &UpstreamError{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(respBody)),
		}
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (c *Client) endpointURL(endpoint string) string {
	return fmt.Sprintf("%s/waInstance%s/%s/%s", c.baseURL, c.idInstance, endpoint, c.apiToken)
}

func (c *Client) validateConfigured() error {
	if strings.TrimSpace(c.baseURL) == "" {
		return ErrMissingBaseURL
	}
	if c.httpClient == nil {
		return ErrMissingHTTPClient
	}
	if strings.TrimSpace(c.idInstance) == "" {
		return ErrMissingIDInstance
	}
	if strings.TrimSpace(c.apiToken) == "" {
		return ErrMissingAPIToken
	}

	return nil
}
