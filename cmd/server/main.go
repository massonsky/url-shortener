package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url_shortener/internal/config"
	"url_shortener/internal/repository/postgres"
	"url_shortener/internal/repository/redis"
	"url_shortener/internal/service"
	"url_shortener/internal/transport/_http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	rdb := redis.NewClient(cfg)
	defer rdb.Close()

	pgRepo := postgres.NewURLRepository(db)
	urlService := service.NewURLService(pgRepo, rdb) // ✅ Передаём напрямую

	router := _http.SetupRouter(urlService)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
