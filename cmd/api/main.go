package main

import (
	"log"
	"net/http"
	"time"

	hs "github.com/ddddddO/tag-mng/internal/api/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("launch api server")

	router := gin.Default()

	// session実装：https://re-engines.com/2020/03/02/go%E3%83%95%E3%83%AC%E3%83%BC%E3%83%A0%E3%83%AF%E3%83%BC%E3%82%AFgin%E3%81%A7%E3%83%9F%E3%83%89%E3%83%AB%E3%82%A6%E3%82%A7%E3%82%A2%E3%82%92%E4%BD%BF%E3%81%A3%E3%81%A6%E3%83%AD%E3%82%B0%E3%82%A4/
	sessionSec := "sessionsecret" // FIXME: Using os.Getenv or crypto/rand
	store := cookie.NewStore([]byte(sessionSec))
	router.Use(sessions.Sessions("tag-mng-session", store))
	router.Use(checkSession())

	// cors実装：https://qiita.com/MasashiFujiike/items/7844150ce75d71a417ad
	router.Use(cors.New(cors.Config{
		// 許可したいHTTPメソッドの一覧
		AllowMethods: []string{
			"GET",
			"OPTIONS",
			"PATCH",
			"POST",
			"DELETE",
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
		// ref: https://developer.mozilla.org/ja/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
		AllowCredentials: true,
		MaxAge:           30 * time.Second,
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
	router.POST("/memodetail", hs.MemoDetailCreateHandler)
	// メモ更新API
	router.PATCH("/memodetail", hs.MemoDetailUpdateHandler)
	// メモ削除API
	router.DELETE("/memodetail", hs.MemoDetailDeleteHandler)
	// タグ一覧返却API
	router.GET("/tags", hs.TagListHandler)
	// タグ詳細返却API
	router.GET("/tagdetail", hs.TagDetailHandler)
	// タグ更新API
	// タグ削除API

	router.Run(":8082")
}

func checkSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("in checkSession")
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		path := c.FullPath()
		if path == "/health" || path == "/auth" {
			c.Next()
			return
		}

		session := sessions.Default(c)
		if session.Get("RANDOM_AUTHED_STRING") != "tmp_authed_token" { // FIXME:
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "required to athenticate",
			})
			c.Abort()
		}

		c.Next()
	}
}
