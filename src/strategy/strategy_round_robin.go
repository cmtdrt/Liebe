package strategy

import (
	"net/http"
	"sync"
)

type roundRobinStrategy struct {
	mu    sync.Mutex
	index int
}

func newRoundRobinStrategy() *roundRobinStrategy {
	return &roundRobinStrategy{
		index: 0,
	}
}

func (s *roundRobinStrategy) Next(healthy []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(healthy) == 0 {
		return "", http.ErrServerClosed
	}
	u := healthy[s.index%len(healthy)]
	s.index = (s.index + 1) % len(healthy)
	return u, nil
}
