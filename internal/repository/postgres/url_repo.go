package postgres

import (
	"context"
	"database/sql"
	"url_shortener/internal/domain"
	"url_shortener/internal/utils"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(ctx context.Context, originalURL string) (*domain.URL, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO urls(original_url) VALUES ($1) RETURNING id`,
		originalURL,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	shortCode := utils.Encode(id)
	_, err = r.db.ExecContext(ctx,
		`UPDATE urls SET short_code = $1 WHERE id = $2`,
		shortCode, id,
	)
	if err != nil {
		return nil, err
	}

	return &domain.URL{
		ID:          id,
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}, nil
}

func (r *URLRepository) IncrementClick(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE urls SET click_count = click_count + 1 WHERE id = $1`, id,
	)
	return err
}
func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	var url domain.URL
	err := r.db.QueryRowContext(ctx,
		`SELECT id, original_url, short_code FROM urls WHERE short_code = $1`,
		shortCode,
	).Scan(&url.ID, &url.OriginalURL, &url.ShortCode)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) GetByOriginalURL(ctx context.Context, originalURL string) (*domain.URL, error) {
	var url domain.URL
	err := r.db.QueryRowContext(ctx,
		`SELECT id, original_url, short_code FROM urls WHERE original_url = $1`,
		originalURL,
	).Scan(&url.ID, &url.OriginalURL, &url.ShortCode)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) Close() error {
	return r.db.Close()
}
