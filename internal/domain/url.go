package domain

import "time"

type URL struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
	ClickCount  int64     `json:"click_count"`
}
