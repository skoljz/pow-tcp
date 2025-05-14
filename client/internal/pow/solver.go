package pow

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
)

type Solver interface {
	Solve(ctx context.Context, challenge []byte) ([]byte, error)
}

type hashcash struct {
	targetSize uint64
}

func NewSolver(targetSize uint64) Solver {
	return &hashcash{targetSize: targetSize}
}

func (s *hashcash) Solve(ctx context.Context, challenge []byte) ([]byte, error) {
	if len(challenge) < int(s.targetSize) {
		return nil, fmt.Errorf("invalid challenge len=%d", len(challenge))
	}

	target := binary.BigEndian.Uint64(challenge[:s.targetSize])
	nonce := make([]byte, s.targetSize)
	data := make([]byte, len(challenge)+int(s.targetSize))
	copy(data, challenge)

	for attempt := uint64(0); attempt < math.MaxUint64; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		binary.BigEndian.PutUint64(nonce, attempt)
		copy(data[len(challenge):], nonce)

		sum := sha256.Sum256(data)
		if binary.BigEndian.Uint64(sum[:s.targetSize]) < target {
			out := make([]byte, s.targetSize)
			copy(out, nonce)
			return out, nil
		}
	}

	return nil, fmt.Errorf("nonce not found")
}
