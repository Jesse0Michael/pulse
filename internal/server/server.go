package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jesse0michael/pulse/internal/service"
)

type Pulser interface {
	Summary(ctx context.Context, req service.SummaryRequest) (string, error)
}

type Config struct {
	Port    int           `envconfig:"PORT" default:"8080"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"30s"`
	Github  service.GithubConfig
	AI      service.OpenAIConfig
}

type Server struct {
	*http.Server
	router *mux.Router
	pulser Pulser
}

func New(cfg Config, pulser Pulser) *Server {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(notAllowed)
	router.Use(handlers.CORS())

	server := &Server{
		Server: &http.Server{
			Handler:     router,
			Addr:        fmt.Sprintf(":%d", cfg.Port),
			ReadTimeout: cfg.Timeout,
		},
		router: router,
		pulser: pulser,
	}

	server.route()

	return server
}

func (s *Server) route() {
	s.router.HandleFunc("/summary/github/users/{username}", s.summary()).Methods("GET").Name("githubUserSummary")
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"error":"page not found"}`))
}

func notAllowed(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"error":"method not allowed"}`))
}

func writeError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err.Error())))
}
