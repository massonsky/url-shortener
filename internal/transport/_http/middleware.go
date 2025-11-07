package _http

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func RateLimit(store map[string]int64, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			now := time.Now().Unix()
			if last, ok := store[ip]; ok && now-last < int64(window.Seconds()) {
				http.Error(w, "rate limited", http.StatusTooManyRequests)
				return
			}
			store[ip] = now
			next.ServeHTTP(w, r)
		})
	}
}
