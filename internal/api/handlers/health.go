package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 1",
		// })
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = conn.Query("SELECT 1")
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": "failed to connect db 2",
		// })
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
