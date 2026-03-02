package strategy

import "net/http"

type roundRobinStrategy struct {
	index int
}

func (s *roundRobinStrategy) Next(healthy []string) (string, error) {
	if len(healthy) == 0 {
		return "", http.ErrServerClosed
	}
	u := healthy[s.index%len(healthy)]
	s.index = (s.index + 1) % len(healthy)
	return u, nil
}

