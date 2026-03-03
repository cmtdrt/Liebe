package strategy

import "liebe/src/config"

// Chooses the next healthy upstream to use.
type StrategyChooser interface {
	Next(healthy []string) (string, error)
}

// Optionally implemented by strategies that need to track per-upstream connection state.
type ConnectionAwareStrategy interface {
	StrategyChooser
	OnRequestStart(upstream string)
	OnRequestEnd(upstream string)
}

// Builds the appropriate strategy implementation from configuration.
func NewStrategyChooser(s config.Strategy) StrategyChooser {
	switch s {
	case config.StrategyRoundRobin:
		return newRoundRobinStrategy()
	case config.StrategyRandom:
		return newRandomStrategy()
	case config.StrategyLeastConnections:
		return newLeastConnectionsStrategy()
	default:
		// Should not happen because config.LoadConfig already validated the value.
		return &roundRobinStrategy{}
	}
}
