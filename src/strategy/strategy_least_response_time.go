package strategy

import (
	"net/http"
	"sync"
	"time"
)

type leastResponseTimeStrategy struct {
	mu               sync.Mutex
	totalDurations   map[string]time.Duration
	requestCounts    map[string]int
}

func newLeastResponseTimeStrategy() *leastResponseTimeStrategy {
	return &leastResponseTimeStrategy{
		totalDurations: make(map[string]time.Duration),
		requestCounts:  make(map[string]int),
	}
}

func (s *leastResponseTimeStrategy) Next(healthy []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(healthy) == 0 {
		return "", http.ErrServerClosed
	}

	// Pick the healthy upstream with the lowest average response time.
	minIdx := 0
	minAvg := s.avg(healthy[0])
	for i := 1; i < len(healthy); i++ {
		u := healthy[i]
		avg := s.avg(u)
		if avg < minAvg {
			minAvg = avg
			minIdx = i
		}
	}
	return healthy[minIdx], nil
}

func (s *leastResponseTimeStrategy) OnRequestComplete(upstream string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalDurations[upstream] += duration
	s.requestCounts[upstream]++
}

func (s *leastResponseTimeStrategy) avg(upstream string) time.Duration {
	total := s.totalDurations[upstream]
	count := s.requestCounts[upstream]
	if count == 0 {
		// If we have no data yet, treat average as 0 so new upstreams are not penalized.
		return 0
	}
	return total / time.Duration(count)
}

