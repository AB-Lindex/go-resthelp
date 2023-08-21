package resthelp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
)

// Response is the response of a request.
type Response struct {
	helper    *Helper
	Header    http.Header
	resp      *http.Response
	respError error
	body      []byte
	bodyError error
}

// Do executes the request and returns the response.
func (req *Request) Do() (*Response, error) {
	intResp, err := req.helper.client.Do(req.req)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		helper:    req.helper,
		resp:      intResp,
		respError: err,
		Header:    intResp.Header,
	}

	return resp, nil
}

// Close closes the response body. (should always be called after Do())
func (r *Response) Close() {
	if r == nil || r.resp == nil {
		return
	}
	if r.resp.Body != nil {
		_ = r.resp.Body.Close()
	}
}

// Status returns the status code of the response.
func (r *Response) Status() int {
	return r.resp.StatusCode
}

// Error returns the error message of the response (or the status since it's "now" considered an error)
func (r *Response) Error() string {
	if r.respError != nil {
		return r.respError.Error()
	}
	return r.resp.Status
}

// IsOK returns true if the response is considered OK.
func (r *Response) IsOK() bool {
	return r.respError != nil ||
		r.resp.StatusCode == http.StatusOK ||
		r.resp.StatusCode == http.StatusNoContent ||
		r.resp.StatusCode == http.StatusCreated
}

func (r *Response) getBody() ([]byte, error) {
	if r.resp.Body != nil {
		r.body, r.bodyError = io.ReadAll(r.resp.Body)
		r.resp.Body = nil
	}
	return r.body, r.bodyError
}

// BodyBytes returns the body of the response as a []byte.
func (r *Response) BodyBytes() ([]byte, error) {
	return r.getBody()
}

const (
	ctJSON = "application/json"
	ctXML  = "application/xml"
)

type errorUnknownCT string

// Error returns a content-type error.
func (e errorUnknownCT) Error() string {
	return fmt.Sprintf("unknown content-type: '%s'", string(e))
}

// ParseContentType returns the simple content-type of the response.
func (r *Response) ParseContentType() string {
	ct := r.Header.Get("Content-Type")
	mediatype, _, _ := mime.ParseMediaType(ct)
	return mediatype
}

// Parse parses the response into the supplied value (using the 'Content-Type' header)
func (r *Response) Parse(v interface{}) error {
	ct := r.ParseContentType()
	if fn, ok := r.helper.parsers[ct]; ok {
		if fn != nil {
			return fn(v)
		}
	}

	switch ct {
	case ctJSON:
		return r.ParseJSON(v)
	case ctXML:
		return r.ParseXML(v)
	}

	return errorUnknownCT(ct)
}

// ParseJSON parses the response into the supplied value as JSON.
func (r *Response) ParseJSON(v interface{}) error {
	data, err := r.getBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ParseXML parses the response into the supplied value as XML.
func (r *Response) ParseXML(v interface{}) error {
	data, err := r.getBody()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}
