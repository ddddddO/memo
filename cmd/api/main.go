package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/ddddddO/tag-mng/api/handler"
	"github.com/ddddddO/tag-mng/api/usecase"
	"github.com/ddddddO/tag-mng/repository/postgres"
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

	// https://tutuz-tech.hatenablog.com/entry/2020/03/24/170159
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(15)
	db.SetConnMaxLifetime(3 * time.Minute)

	store, err := genPostgresStore(db)
	if err != nil {
		log.Fatal(err)
	}
	router.Use(checkSession(store))

	// ヘルスチェック
	healthRepository := postgres.NewHealthRepository(db)
	healthUsecase := usecase.NewHealth(healthRepository)
	healthHandler := handler.NewHealth(healthUsecase)
	router.Get("/health", healthHandler.Check)

	// 認証API
	userRepository := postgres.NewUserRepository(db)
	authUsecase := usecase.NewAuth(userRepository, store)
	authHandler := handler.NewAuth(authUsecase)

	router.Post("/auth", authHandler.Auth)

	memoRepository := postgres.NewMemoRepository(db)
	memoUsecase := usecase.NewMemo(memoRepository)
	memoHandler := handler.NewMemo(memoUsecase)
	router.Route("/memos", func(r chi.Router) {
		// メモ一覧返却API
		r.Get("/", memoHandler.List)
		// メモ新規作成API
		r.Post("/", memoHandler.Create)
		// メモ更新API
		r.Patch("/{id}", memoHandler.Update)
		// メモ削除API
		r.Delete("/{id}", memoHandler.Delete)
		// メモ詳細返却API
		r.Get("/{id}", memoHandler.Detail)
	})

	tagRepository := postgres.NewTagRepository(db)
	tagUsecase := usecase.NewTag(tagRepository)
	tagHandler := handler.NewTag(tagUsecase)
	router.Route("/tags", func(r chi.Router) {
		// タグ一覧返却API
		r.Get("/", tagHandler.List)
		// タグ新規作成API
		r.Post("/", tagHandler.Create)
		// タグ更新API
		r.Patch("/{id}", tagHandler.Update)
		// タグ削除API
		r.Delete("/{id}", tagHandler.Delete)
		// タグ詳細返却API
		r.Get("/{id}", tagHandler.Detail)
	})

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
