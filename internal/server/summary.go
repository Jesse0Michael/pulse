package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jesse0michael/go-request"
	"github.com/jesse0michael/pulse/internal/service"
)

func (s *Server) summary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.SummaryRequest
		if err := request.Decode(r, &req); err != nil {
			slog.With("error", err).Error("failed to decode request body")
			writeError(w, http.StatusBadRequest, err)
			return
		}

		summary, err := s.pulser.Summary(r.Context(), req)
		if err != nil {
			slog.With("error", err).Error("failed to get summary")
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		b, _ := json.Marshal(summary)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}
