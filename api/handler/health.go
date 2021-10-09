package handler

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

type healthUsecase interface {
	Check() error
}

type healthHandler struct {
	usecase healthUsecase
}

func NewHealth(usecase healthUsecase) *healthHandler {
	return &healthHandler{
		usecase: usecase,
	}
}

func (h *healthHandler) Check(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.Check(); err != nil {
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
