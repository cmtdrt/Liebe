package strategy

import (
	"liebe/src/config"
	"time"
)

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

// Optionally implemented by strategies that adjust routing based on per-upstream response times.
type ResponseTimeAwareStrategy interface {
	StrategyChooser
	OnRequestComplete(upstream string, duration time.Duration)
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
	case config.StrategyLeastResponseTime:
		return newLeastResponseTimeStrategy()
	default:
		// Should not happen because config.LoadConfig already validated the value.
		return &roundRobinStrategy{}
	}
}
