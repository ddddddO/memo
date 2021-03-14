package api

import (
	"encoding/json"
	"net/http"
)

func errResponse(w http.ResponseWriter, status int, message string, err error) {
	res := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
