package strategy

import (
	"liebe/src/config"
)

// Chooses the next healthy upstream to use.
type StrategyChooser interface {
	Next(healthy []string) (string, error)
}

// Builds the appropriate strategy implementation from configuration.
func NewStrategyChooser(s config.Strategy) StrategyChooser {
	switch s {
	case config.StrategyRoundRobin:
		return &roundRobinStrategy{}
	case config.StrategyRandom:
		return newRandomStrategy()
	default:
		// Should not happen because config.LoadConfig already validated the value.
		return &roundRobinStrategy{}
	}
}

