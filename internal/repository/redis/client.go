package redis

import (
	"context"
	"time"

	"url_shortener/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

// NewClient создаёт новый Redis-клиент по конфигурации
func NewClient(cfg *config.Config) *Client {
	opts := &redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: 20, // настройка пула под нагрузку
	}

	// Добавим таймауты для отказоустойчивости
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second
	opts.DialTimeout = 5 * time.Second

	client := redis.NewClient(opts)

	// Проверим соединение
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return &Client{client: client}
}

// Get оборачивает вызов Redis GET
func (r *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

// Set оборачивает Redis SET с TTL
func (r *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

// Incr — инкремент счётчика (полезно для статистики)
func (r *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	return r.client.Incr(ctx, key)
}

// Close закрывает соединение (опционально для graceful shutdown)
func (r *Client) Close() error {
	return r.client.Close()
}
