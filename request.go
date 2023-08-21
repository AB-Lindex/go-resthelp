package resthelp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	headerContentType = "Content-Type"
)

// Request is how you create a request to an API.
type Request struct {
	helper  *Helper
	req     *http.Request
	method  string
	url     *url.URL
	body    io.Reader
	headers map[string]string
	query   map[string]string
}

// RequestOption is a function that modifies the Request object.
type RequestOption func(*Request)

// Get creates a GET request to the given path
func (h *Helper) Get(path string, opts ...RequestOption) (*Request, error) {
	return h.NewRequest("GET", path, opts...)
}

// Post creates a POST request to the given path
func (h *Helper) Post(path string, opts ...RequestOption) (*Request, error) {
	return h.NewRequest("POST", path, opts...)
}

// Put creates a PUT request to the given path
func (h *Helper) Put(path string, opts ...RequestOption) (*Request, error) {
	return h.NewRequest("PUT", path, opts...)
}

// Patch creates a PATCH request to the given path
func (h *Helper) Patch(path string, opts ...RequestOption) (*Request, error) {
	return h.NewRequest("PATCH", path, opts...)
}

// Delete creates a DELETE request to the given path
func (h *Helper) Delete(path string, opts ...RequestOption) (*Request, error) {
	return h.NewRequest("DELETE", path, opts...)
}

// NewRequest creates a new Request object with the given options.
func (h *Helper) NewRequest(method, path string, opts ...RequestOption) (*Request, error) {
	var err error
	if h.baseURL != "" {
		path, err = url.JoinPath(h.baseURL, path)
		if err != nil {
			return nil, err
		}
	}
	parsedURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	r := &Request{
		helper: h,
		method: method,
		url:    parsedURL,
		body:   nil,
	}

	// save helper-headers to interim header-map
	for k, v := range h.headers {
		r.AddHeader(k, v)
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.body != nil {
		r.req, err = http.NewRequest(method, r.url.String(), r.body)
	} else {
		r.req, err = http.NewRequest(method, r.url.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// apply interim headers to actual request
	for k, v := range r.headers {
		r.req.Header.Add(k, v)
	}

	return r, nil
}

// Request returns the URL of the request.
func (r *Request) Request() string {
	return r.req.URL.String()
}

// WithQuery adds a query-parameter to the request.
func WithQuery(k, v string) func(*Request) {
	return func(r *Request) {
		r.AddQuery(k, v)
	}
}

// WithRequestHeader adds a header to the request.
func WithRequestHeader(k, v string) func(*Request) {
	return func(r *Request) {
		r.AddHeader(k, v)
	}
}

// WithBody sets the body of the request.
func WithBody(body io.Reader) func(*Request) {
	return func(r *Request) {
		r.body = body
	}
}

// WithJSON sets the body of the request using JSON of the supplied value.
func WithJSON(v interface{}) func(*Request) {
	return func(r *Request) {
		buf, err := json.Marshal(v)
		if err != nil {
			r.body = nil
			return
		}
		r.AddHeader(headerContentType, "application/json")
		r.AddHeader("Content-Length", fmt.Sprint(len(buf)))
		r.body = bytes.NewReader(buf)
	}
}

// WithXML sets the body of the request using XML of the supplied value.
func WithXML(v interface{}) func(*Request) {
	return func(r *Request) {
		buf, err := xml.Marshal(v)
		if err != nil {
			r.body = nil
			return
		}
		r.AddHeader(headerContentType, "application/xml")
		r.AddHeader("Content-Length", fmt.Sprint(len(buf)))
		r.body = bytes.NewReader(buf)
	}
}

// WithContentType sets the Content-Type header of the request.
func WithContentType(ct string) func(*Request) {
	return func(r *Request) {
		r.AddHeader(headerContentType, ct)
	}
}

// Request methods

// AddHeader adds a header to the request.
func (r *Request) AddHeader(k, v string) {
	if r.req != nil {
		r.req.Header.Add(k, v)
		return
	}

	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[k] = v
}

// AddQuery adds a query parameter to the request.
func (r *Request) AddQuery(k, v string) {
	if r.req != nil {
		q := r.req.URL.Query()
		q.Add(k, v)
		r.req.URL.RawQuery = q.Encode()
		return
	}

	if r.query == nil {
		r.query = make(map[string]string)
	}
	r.query[k] = v
}
