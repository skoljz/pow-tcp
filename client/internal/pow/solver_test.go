package pow_test

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"github.com/skoljz/pow_tcp_client/internal/pow"
	"math"
	"testing"
)

func challenge(target uint64, payload []byte) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, target)
	return append(buf, payload...)
}

func TestSolve_Success(t *testing.T) {
	const targetSize = 8
	const target = math.MaxUint64 / 12
	payload := []byte("faraway")

	c := challenge(target, payload)
	solver := pow.NewSolver(targetSize)

	nonce, err := solver.Solve(context.Background(), c)
	if err != nil {
		t.Fatalf("Solve() returned error: %v", err)
	}
	if len(nonce) != targetSize {
		t.Fatalf("nonce length = %d; want %d", len(nonce), targetSize)
	}

	full := append(c, nonce...)
	sum := sha256.Sum256(full)
	hval := binary.BigEndian.Uint64(sum[:8])
	t.Log(hval, target, sum[:8])
	if hval > target {
		t.Fatalf(
			"invalid proof-of-work: hash = %x (uint64=%d) > target=%d",
			sum[:8], hval, target,
		)
	}
}

func TestSolve_InvalidChallenge(t *testing.T) {
	solver := pow.NewSolver(8)
	_, err := solver.Solve(context.Background(), []byte{1, 2, 3, 4})
	if err == nil {
		t.Fatal("expected error")
	}
}
