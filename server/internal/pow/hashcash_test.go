package pow

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestVerify_Success(t *testing.T) {
	p, _ := New(4)
	ch, _ := p.Challenge()
	nonce := solveCh(p, ch)
	if err := p.Verify(ch, nonce); err != nil {
		t.Errorf("Verify failed for valid nonce: %v", err)
	}
}

func TestVerify_Failure(t *testing.T) {
	p, _ := New(6)
	ch, _ := p.Challenge()
	bad := make([]byte, nonceSize) // all zeros
	if err := p.Verify(ch, bad); err == nil {
		t.Error("expected failure for invalid nonce, but got success")
	}
}

func solveCh(provider *HashcashProvider, ch []byte) []byte {
	nonce := make([]byte, nonceSize)

	for attemps := uint64(0); attemps < math.MaxUint64; attemps++ {
		binary.BigEndian.PutUint64(nonce, attemps)
		err := provider.Verify(ch, nonce)

		if err == nil {
			return append([]byte(nil), nonce...)
		}
	}
	return nil
}
