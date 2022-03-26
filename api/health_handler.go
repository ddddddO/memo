package api

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

type healthRepository interface {
	Check() error
}

type healthHandler struct {
	repo healthRepository
}

func NewHealthHandler(repo healthRepository) *healthHandler {
	return &healthHandler{
		repo: repo,
	}
}

func (h *healthHandler) Check(w http.ResponseWriter, r *http.Request) {
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
