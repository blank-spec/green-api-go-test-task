package greenapi

import (
	"errors"
	"fmt"
)

var (
	ErrMissingBaseURL    = errors.New("baseURL is required")
	ErrMissingHTTPClient = errors.New("http client is required")
	ErrMissingIDInstance = errors.New("idInstance is required")
	ErrMissingAPIToken   = errors.New("apiTokenInstance is required")
)

type UpstreamError struct {
	StatusCode int
	Body       string
}

func (e *UpstreamError) Error() string {
	return fmt.Sprintf("green api returned status %d", e.StatusCode)
}
