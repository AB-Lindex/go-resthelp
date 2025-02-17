package resthelp

import (
	"testing"
)

func TestRequest_AddQuery(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
		want   string
	}{
		{"no-param", nil, "http://localhost/test"},
		{"string", map[string]string{"a": "value"}, "http://localhost/test?a=value"},
	}

	base := New(WithBaseURL("http://localhost"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := base.Get("/test")
			for k, v := range tt.fields {
				r.AddQuery(k, v)
			}
			if r.Request() != tt.want {
				t.Errorf("Request.AddQuery() = %v, want %v", r.Request(), tt.want)
			}
		})
	}
}

func TestWithQuery(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
		want   string
	}{
		{"no-param", nil, "http://localhost/test"},
		{"string", map[string]string{"a": "value"}, "http://localhost/test?a=value"},
	}

	base := New(WithBaseURL("http://localhost"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []RequestOption

			for k, v := range tt.fields {
				opts = append(opts, WithQuery(k, v))
			}
			r, _ := base.Get("/test", opts...)
			if r.Request() != tt.want {
				t.Errorf("Request.WithQuery() = %v, want %v", r.Request(), tt.want)
			}
		})
	}
}
