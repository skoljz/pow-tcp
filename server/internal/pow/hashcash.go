package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	challengeSize = 16
	nonceSize     = 8
	minComplexity = 1
	maxComplexity = 24
)

var (
	ErrInvalidComplexity = errors.New("pow: invalid complexity")
	ErrChallengeSize     = errors.New("pow: invalid challenge size")
	ErrSolutionSize      = errors.New("pow: invalid solution size")
	ErrVerification      = errors.New("pow: verification failed")
)

type Provider interface {
	Challenge() ([]byte, error)
	Verify(ch, nonce []byte) error
}

type HashcashProvider struct {
	complexity uint8
}

var ()

func New(complexity uint8) (*HashcashProvider, error) {
	if complexity < minComplexity || complexity > maxComplexity {
		return nil, ErrInvalidComplexity
	}

	return &HashcashProvider{complexity: complexity}, nil
}

func (p *HashcashProvider) Challenge() ([]byte, error) {
	ch := make([]byte, challengeSize)

	bitsToShift := 64 - int(p.complexity)
	var target uint64 = 1
	for i := 0; i < bitsToShift; i++ {
		target *= 2
	}

	targetBytes := ch[:nonceSize]
	binary.BigEndian.PutUint64(targetBytes, target)

	payload := ch[nonceSize:]
	bytesRead, err := rand.Read(payload)
	if err != nil {
		return nil, fmt.Errorf("random bytes: %w", err)
	}

	expected := len(payload)
	if bytesRead != expected {
		return nil, fmt.Errorf("random read: expected %d bytes, got %d", expected, bytesRead)
	}

	return ch, nil
}

func (p *HashcashProvider) Verify(ch, nonce []byte) error {
	if len(ch) != challengeSize {
		return ErrChallengeSize

	}
	if len(nonce) != nonceSize {
		return ErrSolutionSize
	}

	target := binary.BigEndian.Uint64(ch[:nonceSize])
	hash := sha256.Sum256(append(ch, nonce...))

	if binary.BigEndian.Uint64(hash[:nonceSize]) >= target {
		return ErrVerification
	}

	return nil
}
