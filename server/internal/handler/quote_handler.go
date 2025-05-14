package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"

	"log/slog"

	"github.com/skoljz/pow_tcp/internal/config"
	"github.com/skoljz/pow_tcp/internal/pow"
	"github.com/skoljz/pow_tcp/internal/storage/quotes"
	"github.com/skoljz/pow_tcp/internal/utils"
)

const QuoteOfTheDay = "Best Wishes, Reviewer! ~ Semyon Kolenkov"

type QuoteHandler struct {
	cfg      *config.Config
	log      *slog.Logger
	storage  quotes.Storage
	provider pow.Provider
}

func NewQuoteHandler(cfg *config.Config, log *slog.Logger,
	st quotes.Storage, prov pow.Provider) *QuoteHandler {

	return &QuoteHandler{cfg: cfg, log: log, storage: st, provider: prov}
}

func (h *QuoteHandler) Handle(ctx context.Context, c net.Conn) {
	defer c.Close()

	go func() {
		<-ctx.Done()
		_ = c.SetDeadline(time.Now())
	}()

	_ = c.SetDeadline(time.Now().Add(h.cfg.ReadTimeout))

	addr := c.RemoteAddr().String()
	start := time.Now()

	if err := h.challengeResponse(c); err != nil {
		h.log.Warn("session terminated", "client", addr, "error", err)
		return
	}
	if err := h.sendQuote(c); err != nil {
		h.log.Error("send quote", "client", addr, "error", err)
		return
	}

	h.log.Info("done",
		"client", addr,
		"latency", time.Since(start).String(),
	)
}

func (h *QuoteHandler) challengeResponse(c net.Conn) error {
	ch, err := h.provider.Challenge()
	if err != nil {
		return fmt.Errorf("failed to generate challenge: %w", err)
	}
	if err := utils.WriteMsg(c, ch); err != nil {
		return fmt.Errorf("failed to write challenge: %w", err)
	}

	sol, err := utils.ReadMsg(c)
	if err != nil {
		return fmt.Errorf("failed to read solution: %w", err)
	}
	if err := h.provider.Verify(ch, sol); err != nil {
		return fmt.Errorf("verify failed: %w", err)
	}

	h.log.Debug("challenge solved", "challenge", hex.EncodeToString(ch))
	return nil
}

func (h *QuoteHandler) sendQuote(c net.Conn) error {
	quote, err := h.storage.Random()
	if err != nil {
		if errors.Is(err, quotes.ErrNoQuotes) {
			quote = QuoteOfTheDay
		} else {
			return err
		}
	}

	return utils.WriteMsg(c, []byte(quote))
}
