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
			errResponse(w, http.StatusInternalServerError, "failed to connect db 1", err)
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
