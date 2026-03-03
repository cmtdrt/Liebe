package strategy

import (
	"math/rand"
	"net/http"
	"time"
)

type randomStrategy struct {
	rnd *rand.Rand
}

func newRandomStrategy() *randomStrategy {
	return &randomStrategy{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *randomStrategy) Next(healthy []string) (string, error) {
	if len(healthy) == 0 {
		return "", http.ErrServerClosed
	}
	return healthy[s.rnd.Intn(len(healthy))], nil
}
