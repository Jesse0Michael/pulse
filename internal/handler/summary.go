package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jesse0michael/pkg/http/errors"
	"github.com/jesse0michael/pkg/http/server"
	"github.com/jesse0michael/pulse/internal/service"
)

type Summary struct {
	Summary string `json:"summary"`
}

func (s *Handler) summary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		q := r.URL.Query()
		req := service.SummaryRequest{
			Username:     vars["username"],
			Organization: q.Get("organization"),
			Repository:   q.Get("repository"),
		}
		if startDate := q.Get("startDate"); startDate != "" {
			if t, err := time.Parse(time.RFC3339, startDate); err == nil {
				req.StartDate = &t
			}
		}
		if endDate := q.Get("endDate"); endDate != "" {
			if t, err := time.Parse(time.RFC3339, endDate); err == nil {
				req.EndDate = &t
			}
		}

		summary, err := s.pulser.Summary(ctx, req)
		if err != nil {
			slog.ErrorContext(ctx, "failed to get summary", "error", err)
			errors.WriteError(ctx, w, err)
			return
		}

		_ = server.Encode(w, http.StatusOK, &Summary{Summary: summary})
	}
}
