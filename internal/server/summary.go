package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jesse0michael/pulse/internal/service"
)

type Summary struct {
	Summary string `json:"summary"`
}

func (s *Server) summary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		summary, err := s.pulser.Summary(r.Context(), req)
		if err != nil {
			slog.With("error", err).Error("failed to get summary")
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		b, _ := json.Marshal(Summary{Summary: summary})
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}
