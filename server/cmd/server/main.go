package main

import (
	"context"
	"github.com/skoljz/pow_tcp/internal/handler"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/skoljz/pow_tcp/internal/config"
	"github.com/skoljz/pow_tcp/internal/pow"
	"github.com/skoljz/pow_tcp/internal/server"
	"github.com/skoljz/pow_tcp/internal/storage/quotes"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error("config init failed", "error", err)
		os.Exit(1)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel()}))

	st, err := quotes.NewStorage(cfg, log)
	if err != nil {
		log.Error("storage init failed", "error", err)
		os.Exit(1)
	}
	defer st.Close()

	powProv, err := pow.New(cfg.PowComplexity)
	if err != nil {
		log.Error("pow init failed", "error", err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Error("listen failed", "error", err)
		os.Exit(1)
	}

	qHandler := handler.NewQuoteHandler(cfg, log, st, powProv)
	srv := server.New(cfg, log, ln, qHandler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := srv.Run(ctx); err != nil {
		log.Error("server stopped", "error", err)
	}

	log.Info("graceful shutdown complete")
}
