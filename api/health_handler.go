package api

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

type healthUsecase interface {
	Ping() error
}

type healthHandler struct {
	usecase healthUsecase
}

func NewHealthHandler(uc healthUsecase) *healthHandler {
	return &healthHandler{
		usecase: uc,
	}
}

func (h *healthHandler) Check(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.Ping(); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed")
		return
	}

	res := struct {
		Message string `json:"message"`
	}{
		Message: "health ok!",
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed")
		return
	}
}
