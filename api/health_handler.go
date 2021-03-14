package api

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/repository"
)

type healthHandler struct {
	repo repository.HealthRepository
}

func NewHealthHandler(repo repository.HealthRepository) *healthHandler {
	return &healthHandler{
		repo: repo,
	}
}

func (h *healthHandler) Check() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.repo.Check(); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}

		res := struct {
			Message string `json:"message"`
		}{
			Message: "health ok!",
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}
	}
}
