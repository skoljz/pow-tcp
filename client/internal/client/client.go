package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/skoljz/pow_tcp_client/internal/config"
)

const protocol = "tcp"

type Client interface {
	Connect(ctx context.Context) (net.Conn, error)
	RequestChallenge(ctx context.Context, conn net.Conn) ([]byte, error)
	SubmitSolution(ctx context.Context, conn net.Conn, nonce []byte) (string, error)
}

type TCPClient struct {
	cfg *config.Config
}

func New(cfg *config.Config) *TCPClient {
	return &TCPClient{cfg: cfg}
}

func (c *TCPClient) Connect(ctx context.Context) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.DialContext(ctx, protocol, c.cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("connect %s: %w", c.cfg.Addr, err)
	}
	return conn, nil
}

func (c *TCPClient) RequestChallenge(_ context.Context, conn net.Conn) ([]byte, error) {
	var length uint64
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("challenge length: %w", err)
	}
	ch := make([]byte, length)
	if _, err := io.ReadFull(conn, ch); err != nil {
		return nil, fmt.Errorf("challenge data: %w", err)
	}
	return ch, nil
}

func (c *TCPClient) SubmitSolution(_ context.Context, conn net.Conn, nonce []byte) (string, error) {
	if err := binary.Write(conn, binary.BigEndian, uint64(len(nonce))); err != nil {
		return "", fmt.Errorf("nonce length: %w", err)
	}
	if _, err := conn.Write(nonce); err != nil {
		return "", fmt.Errorf("nonce data: %w", err)
	}

	var qlen uint64
	if err := binary.Read(conn, binary.BigEndian, &qlen); err != nil {
		return "", fmt.Errorf("quote length: %w", err)
	}
	buf := make([]byte, qlen)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return "", fmt.Errorf("quote data: %w", err)
	}
	return string(buf), nil
}
