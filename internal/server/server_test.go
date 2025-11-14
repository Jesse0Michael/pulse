package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	server := New(Config{}, nil)

	assert.NotNil(t, server.router, "router should not be nil")
}

func TestServer_route(t *testing.T) {
	server := New(Config{}, nil)

	expected := []string{"githubUserSummary"}
	received := []string{}

	_ = server.router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		received = append(received, route.GetName())
		return nil
	})

	assert.Equal(t, expected, received)
}

func Test_notFound(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	notFound(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.JSONEq(t, `{"error":"page not found"}`, resp.Body.String())
}

func Test_notAllowed(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	notAllowed(resp, req)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
	assert.JSONEq(t, `{"error":"method not allowed"}`, resp.Body.String())
}
