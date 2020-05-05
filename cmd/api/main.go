package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	hs "github.com/ddddddO/tag-mng/internal/api/handlers"
	in "github.com/ddddddO/tag-mng/internal/api/infra"
	uc "github.com/ddddddO/tag-mng/internal/api/usecase"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	log.Println("launch api server")

	router := chi.NewRouter()

	sessionSec := "sessionsecret" // FIXME: Using os.Getenv or crypto/rand
	store := sessions.NewCookieStore([]byte(sessionSec))

	router.Use(checkSession(store))

	// cors: https://github.com/rs/cors#parameters
	c := cors.New(cors.Options{
		AllowedMethods: []string{
			"GET",
			"OPTIONS",
			"PATCH",
			"POST",
			"DELETE",
		},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Origin",
		},
		AllowedOrigins: []string{
			"http://localhost:8080",
			"http://127.0.0.1:8887", // Web Server for Chrome
			"https://app-dot-tag-mng-243823.appspot.com",
		},
		// ref: https://developer.mozilla.org/ja/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
		AllowCredentials: true,
		//MaxAge:           30,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	// ref: chi use cors
	//     https://github.com/rs/cors/blob/master/examples/chi/server.go
	router.Use(c.Handler)

	// TODO: 共通化
	DBDSN := os.Getenv("DBDSN")
	if len(DBDSN) == 0 {
		log.Println("set default DSN")
		DBDSN = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}

	DB, err := sql.Open("postgres", DBDSN)
	if err != nil {
		panic(err)
	}

	user := in.NewUser(DB)
	userUseCase := uc.NewUserUseCase(user)
	authHandler := hs.NewAuthHandler(userUseCase)

	// health
	router.Get("/health", hs.HealthHandler(DB))
	// 認証API
	router.Post("/auth", authHandler.Login(store).(http.HandlerFunc))
	// メモ一覧返却API
	router.Get("/memos", hs.MemoListHandler(DB))
	// メモ詳細返却API
	router.Get("/memodetail", hs.MemoDetailHandler(DB))
	// メモ新規作成API
	router.Post("/memodetail", hs.MemoDetailCreateHandler(DB))
	// メモ更新API
	router.Patch("/memodetail", hs.MemoDetailUpdateHandler(DB))
	// メモ削除API
	router.Delete("/memodetail", hs.MemoDetailDeleteHandler(DB))
	// タグ一覧返却API
	router.Get("/tags", hs.TagListHandler(DB))
	// タグ詳細返却API
	router.Get("/tagdetail", hs.TagDetailHandler(DB))
	// タグ新規作成API
	router.Post("/tagdetail", hs.TagDetailCreateHandler(DB))
	// タグ更新API
	router.Patch("/tagdetail", hs.TagDetailUpdateHandler(DB))
	// タグ削除API
	router.Delete("/tagdetail", hs.TagDetailDeleteHandler(DB))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

// ref: chi middleware
//      https://github.com/go-chi/chi#middleware-handlers
func checkSession(store sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("in checkSession")
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			path := r.URL.EscapedPath()
			log.Println(path)
			if path == "/health" || path == "/auth" {
				next.ServeHTTP(w, r)
				return
			}

			session, _ := store.Get(r, "STORE")
			val, ok := session.Values["authed"].(bool)
			if !ok || !val {
				log.Println("UnAuthenticated")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
