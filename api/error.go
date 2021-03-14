package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func errResponse(w http.ResponseWriter, status int, message string, err error) {
	log.Println(err)

	res := response{
		Message: message,
	}
	resJSON, _ := json.Marshal(res)
	w.WriteHeader(status)
	w.Write([]byte(resJSON))
}
