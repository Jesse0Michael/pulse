package handler

import (
	"context"
	"net/http"

	"github.com/jesse0michael/pulse/internal/service"
)

type Pulser interface {
	Summary(ctx context.Context, req service.SummaryRequest) (string, error)
}

type Handler struct {
	pulser Pulser
}

func New(pulser Pulser) *Handler {
	return &Handler{pulser: pulser}
}

func (s *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /summary/github/users/{username}", s.summary())

	return mux
}
