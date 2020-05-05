package main

import (
	"log"
	"net/http"

	hs "github.com/ddddddO/tag-mng/internal/api/handlers"
	in "github.com/ddddddO/tag-mng/internal/api/infra"
	uc "github.com/ddddddO/tag-mng/internal/api/usecase"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
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

	user := in.NewUser()
	userUseCase := uc.NewUserUseCase(user)
	authHandler := hs.NewAuthHandler(userUseCase)

	// health
	router.Get("/health", hs.HealthHandler)
	// 認証API
	router.Post("/auth", authHandler.Login(store).(http.HandlerFunc))
	// メモ一覧返却API
	router.Get("/memos", hs.MemoListHandler)
	// メモ詳細返却API
	router.Get("/memodetail", hs.MemoDetailHandler)
	// メモ新規作成API
	router.Post("/memodetail", hs.MemoDetailCreateHandler)
	// メモ更新API
	router.Patch("/memodetail", hs.MemoDetailUpdateHandler)
	// メモ削除API
	router.Delete("/memodetail", hs.MemoDetailDeleteHandler)
	// タグ一覧返却API
	router.Get("/tags", hs.TagListHandler)
	// タグ詳細返却API
	router.Get("/tagdetail", hs.TagDetailHandler)
	// タグ新規作成API
	router.Post("/tagdetail", hs.TagDetailCreateHandler)
	// タグ更新API
	router.Patch("/tagdetail", hs.TagDetailUpdateHandler)
	// タグ削除API
	router.Delete("/tagdetail", hs.TagDetailDeleteHandler)

	http.ListenAndServe(":8082", router)
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
