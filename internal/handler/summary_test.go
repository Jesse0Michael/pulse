package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jesse0michael/pulse/internal/service"
	"github.com/stretchr/testify/assert"
)

type MockPulser struct {
	expected service.SummaryRequest
	summary  string
	err      error
}

func (m *MockPulser) Summary(_ context.Context, req service.SummaryRequest) (string, error) {
	if !reflect.DeepEqual(req, m.expected) {
		return "", errors.New("unexpected req")
	}
	return m.summary, m.err
}

func TestServer_summary(t *testing.T) {
	startDate, _ := time.Parse(time.RFC3339, "2023-09-12T00:00:01Z")
	endDate, _ := time.Parse(time.RFC3339, "2023-09-13T23:59:59Z")
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
			name: "successful summary: with params",
			req: httptest.NewRequest(
				http.MethodGet,
				"/summary/github/users/jesse0michael?organization=Jesse0Michael&repository=pulse&startDate=2023-09-12T00:00:01Z&endDate=2023-09-13T23:59:59Z",
				nil,
			),
			pulser: &MockPulser{
				expected: service.SummaryRequest{
					Username:     "jesse0michael",
					Organization: "Jesse0Michael",
					Repository:   "pulse",
					StartDate:    &startDate,
					EndDate:      &endDate,
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
			wantBody: `{"errors":[{"message":"Internal Server Error"}]}`,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.pulser)

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
