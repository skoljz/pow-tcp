package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"log/slog"

	"github.com/skoljz/pow_tcp/internal/config"
)

type ConnHandler interface {
	Handle(ctx context.Context, c net.Conn)
}

type Server struct {
	cfg     *config.Config
	log     *slog.Logger
	ln      net.Listener
	handler ConnHandler

	wg sync.WaitGroup
}

func New(cfg *config.Config, log *slog.Logger,
	ln net.Listener, h ConnHandler) *Server {

	return &Server{cfg: cfg, log: log, ln: ln, handler: h}
}

func (s *Server) Run(ctx context.Context) error {
	defer s.ln.Close()

	s.log.Info("server started", "addr", s.ln.Addr().String())

	go func() {
		<-ctx.Done()
		_ = s.ln.Close()
	}()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				break
			}
			s.log.Warn("failed accept", "error", err)
			continue
		}

		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			s.handler.Handle(ctx, c)
		}(conn)
	}

	return s.waitShutdown()
}

func (s *Server) waitShutdown() error {
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.log.Info("all connections finishd, shutdown complete")
		return nil
	case <-time.After(s.cfg.ShutdownTimeout):
		s.log.Warn("shutdown timeout exceeded", "timeout", s.cfg.ShutdownTimeout)
		return fmt.Errorf("shutdown timed out after %s", s.cfg.ShutdownTimeout)
	}
}
