package api

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// ref: chi middleware
//      https://github.com/go-chi/chi#middleware-handlers
func CheckSession(store sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
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
		}
		return http.HandlerFunc(fn)
	}
}
