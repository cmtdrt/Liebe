package strategy

import (
	"net/http"
	"sync"
)

type leastConnectionsStrategy struct {
	mu                sync.Mutex
	activeConnections map[string]int
}

func newLeastConnectionsStrategy() *leastConnectionsStrategy {
	return &leastConnectionsStrategy{
		activeConnections: make(map[string]int),
	}
}

func (s *leastConnectionsStrategy) Next(healthy []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(healthy) == 0 {
		return "", http.ErrServerClosed
	}

	// Pick the healthy upstream with the fewest active connections.
	minIdx := 0
	minVal := s.activeConnections[healthy[0]]
	for i := 1; i < len(healthy); i++ {
		u := healthy[i]
		if s.activeConnections[u] < minVal {
			minVal = s.activeConnections[u]
			minIdx = i
		}
	}
	return healthy[minIdx], nil
}

func (s *leastConnectionsStrategy) OnRequestStart(upstream string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeConnections[upstream]++
}

func (s *leastConnectionsStrategy) OnRequestEnd(upstream string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeConnections[upstream] > 0 {
		s.activeConnections[upstream]--
	}
}

