package auth_test

import (
	"net/http"
	"testing"

	"github.com/arjunsingh14/chirpy/internal/auth"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantError bool
	}{
		{
			name:      "valid bearer token",
			header:    "Bearer header.payload.signature",
			wantToken: "header.payload.signature",
		},
		{
			name:      "valid bearer token with extra whitespace",
			header:    "  Bearer   header.payload.signature  ",
			wantToken: "header.payload.signature",
		},
		{
			name:      "case-insensitive bearer scheme",
			header:    "bearer header.payload.signature",
			wantToken: "header.payload.signature",
		},
		{
			name:      "missing header",
			wantError: true,
		},
		{
			name:      "missing token",
			header:    "Bearer",
			wantError: true,
		},
		{
			name:      "wrong authorization scheme",
			header:    "Basic credentials",
			wantError: true,
		},
		{
			name:      "too many fields",
			header:    "Bearer token extra",
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			headers := http.Header{}
			if test.header != "" {
				headers.Set("Authorization", test.header)
			}

			token, err := auth.GetBearerToken(headers)
			if test.wantError {
				if err == nil {
					t.Fatalf("GetBearerToken() error = nil, want an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetBearerToken() returned an unexpected error: %v", err)
			}
			if token != test.wantToken {
				t.Errorf("GetBearerToken() token = %q, want %q", token, test.wantToken)
			}
		})
	}
}
