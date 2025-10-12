package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Logger is a middleware that logs HTTP requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Recoverer is a middleware that recovers from panics
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// actorUserIDKey is the context key under which the current actor's user ID is stored.
type actorUserIDKeyType string

const actorUserIDKey actorUserIDKeyType = "actorUserID"

// InjectActorUser returns middleware that injects a fixed actor user ID into the request context.
// This is a placeholder for future authentication integration.
func InjectActorUser(userID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), actorUserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ActorUserIDFromContext retrieves the injected actor user ID from context.
func ActorUserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(actorUserIDKey)
	if v == nil {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}
