package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"test-task/internal/config"
	"test-task/internal/greenapi"
	"test-task/internal/httpapi"
)

const shutdownTimeout = 10 * time.Second

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	app := httpapi.NewApp(greenapi.NewClient(cfg))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErrCh := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", cfg.HTTPAddr)
		serverErrCh <- app.Listen(cfg.HTTPAddr)
	}()

	select {
	case err := <-serverErrCh:
		if err != nil && !errors.Is(err, net.ErrClosed) {
			return fmt.Errorf("listen: %w", err)
		}
		return nil
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	if err := <-serverErrCh; err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("listen: %w", err)
	}

	return nil
}
