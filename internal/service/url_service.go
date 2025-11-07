package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"
	"url_shortener/internal/domain"
	"url_shortener/internal/repository/postgres"
	"url_shortener/internal/repository/redis"
)

var (
	ErrInvalidURL = errors.New("invalid URL")
	ErrNotFound   = errors.New("short URL not found")
)

type URLService struct {
	pgRepo *postgres.URLRepository
	redis  *redis.Client
}

func NewURLService(pg *postgres.URLRepository, redis *redis.Client) *URLService {
	return &URLService{pgRepo: pg, redis: redis}
}

func (s *URLService) Shorten(ctx context.Context, originalURL string) (*domain.URL, error) {
	// Нормализуем URL: убираем пробелы
	originalURL = strings.TrimSpace(originalURL)

	if originalURL == "" {
		return nil, ErrInvalidURL
	}

	// Парсим с использованием url.Parse
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return nil, ErrInvalidURL
	}

	// Проверяем, что есть хост (scheme может быть добавлен handler'ом)
	if parsedURL.Host == "" {
		return nil, ErrInvalidURL
	}

	return s.pgRepo.Create(ctx, originalURL)
}

func (s *URLService) Resolve(ctx context.Context, shortCode string) (string, error) {
	// Сначала проверим кэш
	cached, err := s.redis.Get(ctx, shortCode).Result()
	if err == nil {
		// Фоном обновляем счётчик (без блокирования)
		go func() {
			// Получаем ID из БД для обновления счётчика
			dbURL, _ := s.pgRepo.GetByShortCode(context.Background(), shortCode)
			if dbURL != nil {
				s.pgRepo.IncrementClick(context.Background(), dbURL.ID)
			}
		}()
		return cached, nil
	}

	// Иначе — из БД
	dbURL, err := s.pgRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", ErrNotFound
	}

	// Кэшируем
	s.redis.Set(ctx, shortCode, dbURL.OriginalURL, time.Hour*24)

	// Инкрементим счётчик
	go s.pgRepo.IncrementClick(ctx, dbURL.ID)

	return dbURL.OriginalURL, nil
}
