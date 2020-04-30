package handlers

import (
	"encoding/json"
	"net/http"
)
type response struct {
	Message string `json:"message"`
}

func errResponse(w http.ResponseWriter, status int, message string) {
	res := response{
		Message: message,
	}
	resJSON, _ := json.Marshal(res)
	w.WriteHeader(status)
	w.Write([]byte(resJSON))
}