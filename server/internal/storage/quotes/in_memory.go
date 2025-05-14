package quotes

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type InMemoryStorage struct {
	mu     sync.Mutex
	rand   *rand.Rand
	quotes []string
}

func NewInMemory(path string) (*InMemoryStorage, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	var quotes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		quotes = append(quotes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return &InMemoryStorage{
		quotes: quotes,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (s *InMemoryStorage) Random() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.quotes) == 0 {
		return "", ErrNoQuotes
	}
	return s.quotes[s.rand.Intn(len(s.quotes))], nil
}

func (s *InMemoryStorage) Close() error { return nil }
