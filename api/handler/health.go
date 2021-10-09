package handler

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/repository"
)

type healthHandler struct {
	healthRepo repository.HealthRepository
}

func NewHealth(healthRepo repository.HealthRepository) *healthHandler {
	return &healthHandler{
		healthRepo: healthRepo,
	}
}

func (h *healthHandler) Check(w http.ResponseWriter, r *http.Request) {
	if err := h.healthRepo.Check(); err != nil {
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
