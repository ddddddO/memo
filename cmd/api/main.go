package main

import (
	"log"

	hs "github.com/ddddddO/tag-mng/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("launch api server")

	router := gin.Default()

	router.GET("/health", hs.HealthHandler)

	router.Run(":8082")
}
