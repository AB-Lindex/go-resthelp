package resthelp

import (
	"encoding/base64"
	"net/http"
	"time"
)

const (
	defaultTimeout = 15 * time.Second
)

// Helper is an object that contains the base URL, headers, and client
// and help you create requests to an API.
type Helper struct {
	baseURL string
	headers map[string]string
	client  http.Client
	parsers map[string]Parser
}

// HelperOption is a function that modifies the Helper object.
type HelperOption func(*Helper)

// Parser is a function that parses the response body into an object.
type Parser func(interface{}) error

// New creates a new Helper object with the given options.
func New(opts ...HelperOption) *Helper {
	h := &Helper{
		client: http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// WithBaseURL sets the base URL of the Helper object.
func WithBaseURL(base string) func(*Helper) {
	return func(h *Helper) {
		h.baseURL = base
	}
}

// WithTimeout sets the timeout of http.Client.
func WithTimeout(d time.Duration) func(*Helper) {
	return func(h *Helper) {
		h.client.Timeout = d
	}
}

// WithHeader sets the default header of the Helper object (like same Authorization header)
func WithHeader(k, v string) func(*Helper) {
	return func(h *Helper) {
		if h.headers == nil {
			h.headers = make(map[string]string)
		}
		h.headers[k] = v
	}
}

// WithParser adds a parser to the Helper object for the given MIME type.
func WithParser(mimeType string, fn Parser) func(*Helper) {
	return func(h *Helper) {
		if h.parsers == nil {
			h.parsers = make(map[string]Parser)
		}
		h.parsers[mimeType] = fn
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// WithBasicAuth sets the Authorization header with the given username and password.
func WithBasicAuth(username, password string) func(*Helper) {
	return WithHeader("Authorization", "Basic "+basicAuth(username, password))
}
