package main

import (
	"log"
	"time"

	hs "github.com/ddddddO/tag-mng/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("launch api server")

	router := gin.Default()

	// cors実装：https://qiita.com/MasashiFujiike/items/7844150ce75d71a417ad
	router.Use(cors.New(cors.Config{
		// 許可したいHTTPメソッドの一覧
		AllowMethods: []string{
			"GET",
			"OPTIONS",
			"PATCH",
		},
		// 許可したいHTTPリクエストヘッダの一覧
		AllowHeaders: []string{
			"Accept",
			"Content-Type",
		},
		// 許可したいアクセス元の一覧
		AllowOrigins: []string{
			"http://localhost:8080",
			"http://127.0.0.1:8887", // Web Server for Chrome
		},
		MaxAge: 30 * time.Second,
	}))

	router.GET("/health", hs.HealthHandler)
	// NOTE: tag-mng/web/app.rbから必要なAPIを列挙
	// 認証API
	router.POST("/auth", hs.AuthHandler)
	// メモ一覧返却API
	router.GET("/memos", hs.MemoListHandler)
	// メモ詳細返却API
	router.GET("/memodetail", hs.MemoDetailHandler)
	// メモ新規作成API
	// メモ更新API
	router.PATCH("/memodetail", hs.MemoDetailUpdateHandler)
	// メモ削除API
	// タグ一覧返却API
	// タグ更新API
	// タグ削除API
	router.Run(":8082")
}
