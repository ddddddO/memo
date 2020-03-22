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
	// NOTE: tag-mng/web/app.rbから必要なAPIを列挙
	// 認証API
	router.POST("/auth", hs.AuthHandler)
	// メモ一覧返却API
	// メモ詳細返却API
	// メモ新規作成API
	// メモ更新API
	// メモ削除API
	// タグ一覧返却API
	// タグ更新API
	// タグ削除API
	router.Run(":8082")
}
