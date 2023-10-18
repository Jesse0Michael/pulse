package server

import (
	"context"
	"errors"
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
	summary  string
	err      error
}

func (m *MockPulser) Summary(ctx context.Context, req service.SummaryRequest) (string, error) {
	if req != m.expected {
		return "", fmt.Errorf("unexpected req")
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
	}{
		{
			name: "successful summary",
			req:  httptest.NewRequest(http.MethodGet, "/summary/github/users/jesse0michael", nil),
			pulser: &MockPulser{
				expected: service.SummaryRequest{
					Username: "jesse0michael",
				},
				summary: "Overall, the user jesse0michael has been actively working on multiple repositories",
			},
			wantCode: http.StatusOK,
			wantBody: `{"summary":"Overall, the user jesse0michael has been actively working on multiple repositories"}`,
		},
		{
			name: "failed summary",
			req:  httptest.NewRequest(http.MethodGet, "/summary/github/users/jesse0michael", nil),
			pulser: &MockPulser{
				expected: service.SummaryRequest{
					Username: "jesse0michael",
				},
				summary: "",
				err:     errors.New("test-error"),
			},
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"test-error"}`,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			s := New(Config{}, tt.pulser)

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/summary/github/users/{username}", s.summary())
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
