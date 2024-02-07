package common

import (
	"net/http"
	"time"
)

// NewHTTPClient returns a http client sensibly configured for testing
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
	}
}
