package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"test-task/internal/config"
	"test-task/internal/greenapi"
	"test-task/internal/httpapi"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client := greenapi.NewClient(cfg)
	app := httpapi.NewRouter(client)

	serverErrCh := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", cfg.HTTPAddr)
		serverErrCh <- app.Listen(cfg.HTTPAddr)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrCh:
		if err != nil {
			log.Fatalf("listen: %v", err)
		}
		return
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}

	if err := <-serverErrCh; err != nil {
		log.Fatalf("listen: %v", err)
	}
}
