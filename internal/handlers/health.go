package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func HealthHandler(c *gin.Context)  {
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	conn, err := sql.Open("postgres", DBDSN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 1",
		})
		return
	}

	_, err = conn.Query("SELECT 1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to connect db 2",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "health ok!",
	})
}