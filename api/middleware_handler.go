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
			log.Println("in check session")

			if r.Method == http.MethodOptions {
				log.Println("in options")
				next.ServeHTTP(w, r)
				return
			}

			path := r.URL.EscapedPath()
			if path == "/health" || path == "/auth" {
				log.Println("health or auth path")
				next.ServeHTTP(w, r)
				return
			}

			session, _ := store.Get(r, "STORE")
			val, ok := session.Values["authed"].(bool)
			if !ok || !val {
				log.Println("unauthenticated")

				errResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
