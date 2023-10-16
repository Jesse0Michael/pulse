package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jesse0michael/pulse/internal/service"
	"github.com/stretchr/testify/assert"
)

type MockPulser struct {
	expected service.SummaryRequest
	summary  *service.Summary
	err      error
}

func (m *MockPulser) Summary(ctx context.Context, req service.SummaryRequest) (*service.Summary, error) {
	if req != m.expected {
		return nil, fmt.Errorf("unexpected req")
	}
	return m.summary, m.err
}

func TestServer_summary(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		pulser   Pulser
		wantCode int
		wantBody string
	}{}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			s := New(Config{}, tt.pulser)

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/summary", s.summary())
			router.ServeHTTP(resp, tt.req)

			result := resp.Result()
			assert.Equal(t, tt.wantCode, result.StatusCode)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, resp.Body.String())
			} else {
				assert.Empty(t, resp.Body.String())
			}
			result.Body.Close()
		})
	}
}
