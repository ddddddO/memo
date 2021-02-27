package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

func HealthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := db.Query("SELECT 1")
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed to connect db 2", err)
			return
		}

		type response struct {
			Message string `json:"message"`
		}
		res := response{
			Message: "health ok!",
		}

		resJson, err := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resJson))
	}
}
