package greenapi

import (
	"context"
	"net/http"
	"strings"

	"test-task/internal/config"
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
