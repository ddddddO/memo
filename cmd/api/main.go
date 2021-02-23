package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ddddddO/tag-mng/api"

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
	router.Get("/health", api.HealthHandler(db))
	// 認証API
	router.Post("/auth", api.NewAuthHandler(db, store).(http.HandlerFunc))

	// TODO: /memos
	//        /memos/{id} な形にする
	// メモ一覧返却API
	router.Get("/memos", api.MemoListHandler(db))
	// メモ詳細返却API
	router.Get("/memodetail", api.MemoDetailHandler(db))
	// メモ新規作成API
	router.Post("/memodetail", api.MemoDetailCreateHandler(db))
	// メモ更新API
	router.Patch("/memodetail", api.MemoDetailUpdateHandler(db))
	// メモ削除API
	router.Delete("/memodetail", api.MemoDetailDeleteHandler(db))

	// TODO: /tags
	//        /tags/{id} な形にする
	// タグ一覧返却API
	router.Get("/tags", api.TagListHandler(db))
	// タグ詳細返却API
	router.Get("/tagdetail", api.TagDetailHandler(db))
	// タグ新規作成API
	router.Post("/tagdetail", api.TagDetailCreateHandler(db))
	// タグ更新API
	router.Patch("/tagdetail", api.TagDetailUpdateHandler(db))
	// タグ削除API
	router.Delete("/tagdetail", api.TagDetailDeleteHandler(db))

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

	// ローカル開発時
	if os.Getenv("DEBUG") != "" {
		pgStore.Options = &sessions.Options{
			MaxAge:   60 * 60 * 24,
			Secure:   false,
			SameSite: http.SameSiteDefaultMode,
		}
	} else { // 本番
		pgStore.Options = &sessions.Options{
			MaxAge:   60 * 60 * 6, // Cookieの有効期限。一旦6時間
			Secure:   true,        // trueの時、httpsのみでCookie使用可能
			SameSite: http.SameSiteNoneMode,
		}
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
