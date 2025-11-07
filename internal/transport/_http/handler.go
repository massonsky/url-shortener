package _http

import (
	"encoding/json"
	"net/http"
	"strings"
	"url_shortener/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	urlService *service.URLService
}

func NewHandler(us *service.URLService) *Handler {
	return &Handler{urlService: us}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	cleanURL := strings.TrimSpace(req.URL)
	if cleanURL == "" {
		http.Error(w, `{"error": "URL is required"}`, http.StatusBadRequest)
		return
	}

	// Добавляем схему если её нет
	if !strings.HasPrefix(cleanURL, "http://") && !strings.HasPrefix(cleanURL, "https://") {
		cleanURL = "https://" + cleanURL
	}

	urlObj, err := h.urlService.Shorten(r.Context(), cleanURL)
	if err != nil {
		// Логируем ошибку для отладки
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid URL: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": "http://localhost:8080/" + urlObj.ShortCode,
	})
}
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code") // ← правильно для Chi

	if shortCode == "" {
		http.NotFound(w, r)
		return
	}

	original, err := h.urlService.Resolve(r.Context(), shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}
