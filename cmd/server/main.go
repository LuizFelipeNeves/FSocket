package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fsocket/internal/config"
	"fsocket/internal/handler"
	"fsocket/internal/hub"
	"fsocket/internal/middleware"
)

func main() {
	cfg := config.Load()

	h := hub.New()
	go h.Run()

	mux := http.NewServeMux()

	mux.HandleFunc("/sse", handler.SSE(h))
	mux.HandleFunc("/publish", middleware.Auth(cfg)(handler.Publish(h)))
	mux.HandleFunc("/publish/broadcast", middleware.Auth(cfg)(handler.Broadcast(h)))
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/stats", handler.Stats(h))

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("SSE Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
