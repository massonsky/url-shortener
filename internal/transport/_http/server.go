package _http

import (
	"net/http"
	"url_shortener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRouter(urlService *service.URLService) http.Handler {
	r := chi.NewRouter()
	r.Use(Logging)
	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))

	handler := NewHandler(urlService)
	r.Post("/shorten", handler.Shorten)
	r.Get("/{code}", handler.Redirect)

	return r
}
