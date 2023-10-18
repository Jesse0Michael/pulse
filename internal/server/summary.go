package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jesse0michael/pulse/internal/service"
)

type Summary struct {
	Summary string `json:"summary"`
}

func (s *Server) summary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		req := service.SummaryRequest{
			Username: vars["username"],
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
