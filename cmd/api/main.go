package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	hs "github.com/ddddddO/tag-mng/api/handlers"

	"github.com/antonlindstrom/pgstore"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	log.Println("launch api server")

	router := chi.NewRouter()

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

	db, err := genDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store, err := genPostgresStore(db)
	if err != nil {
		log.Fatal(err)
	}

	router.Use(checkSession(store))

	// health
	router.Get("/health", hs.HealthHandler(db))
	// 認証API
	router.Post("/auth", hs.NewAuthHandler(db, store).(http.HandlerFunc))

	// TODO: /memos
	//        /memos/{id} な形にする
	// メモ一覧返却API
	router.Get("/memos", hs.MemoListHandler(db))
	// メモ詳細返却API
	router.Get("/memodetail", hs.MemoDetailHandler(db))
	// メモ新規作成API
	router.Post("/memodetail", hs.MemoDetailCreateHandler(db))
	// メモ更新API
	router.Patch("/memodetail", hs.MemoDetailUpdateHandler(db))
	// メモ削除API
	router.Delete("/memodetail", hs.MemoDetailDeleteHandler(db))

	// TODO: /tags
	//        /tags/{id} な形にする
	// タグ一覧返却API
	router.Get("/tags", hs.TagListHandler(db))
	// タグ詳細返却API
	router.Get("/tagdetail", hs.TagDetailHandler(db))
	// タグ新規作成API
	router.Post("/tagdetail", hs.TagDetailCreateHandler(db))
	// タグ更新API
	router.Patch("/tagdetail", hs.TagDetailUpdateHandler(db))
	// タグ削除API
	router.Delete("/tagdetail", hs.TagDetailDeleteHandler(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}

func genDB() (*sql.DB, error) {
	dsn := os.Getenv("DBDSN")
	if len(dsn) == 0 {
		log.Println("set default DSN")
		dsn = "host=localhost dbname=tag-mng user=postgres password=postgres sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// old session impl
func genCookieStore() sessions.Store {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "sessionsecret"
	}
	cookieStore := sessions.NewCookieStore([]byte(sessionKey))
	cookieStore.Options = &sessions.Options{
		MaxAge: 60 * 60 * 6, // Cookieの有効期限。一旦6時間
	}
	return cookieStore
}

// new session impl
func genPostgresStore(db *sql.DB) (sessions.Store, error) {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "sessionsecret"
	}

	pgStore, err := pgstore.NewPGStoreFromPool(db, []byte(sessionKey))
	if err != nil {
		return nil, err
	}

	pgStore.Options = &sessions.Options{
		MaxAge:   60 * 60 * 6, // Cookieの有効期限。一旦6時間
		Secure:   true,        // httpsのみ使用可能
		SameSite: http.SameSiteNoneMode,
	}
	return pgStore, nil
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
