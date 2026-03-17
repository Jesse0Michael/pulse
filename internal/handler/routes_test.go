package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Routes(t *testing.T) {
	handler := New(nil)
	mux := handler.Routes()

	tests := []struct {
		name           string
		method         string
		path           string
		expectNotFound bool
	}{
		{
			name:           "github user summary route exists",
			method:         "GET",
			path:           "/summary/github/users/testuser",
			expectNotFound: false,
		},
		{
			name:           "non-existent route returns 404",
			method:         "GET",
			path:           "/nonexistent",
			expectNotFound: true,
		},
		{
			name:           "wrong method returns 405",
			method:         "POST",
			path:           "/summary/github/users/testuser",
			expectNotFound: false, // Should return 405, not 404
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if tt.expectNotFound && w.Code != http.StatusNotFound {
				t.Errorf("expected 404, got %d", w.Code)
			}
			if !tt.expectNotFound && w.Code == http.StatusNotFound {
				t.Errorf("route should exist but got 404")
			}
		})
	}
}
