package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Strategy string

const (
	StrategyRoundRobin      Strategy = "round_robin"
	StrategyRandom          Strategy = "random"
	StrategyLeastConnections Strategy = "least_connections"
)

type Config struct {
	Upstreams   []string    `json:"upstreams"`
	HealthCheck HealthCheck `json:"health_check"`
	Strategy    Strategy    `json:"strategy"`
}

type HealthCheck struct {
	Path     string `json:"path"`
	Interval int    `json:"interval,string"`
	Timeout  int    `json:"timeout,string"`
}

func LoadConfig(path string) (*Config, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}

	if err := validateStrategy(config.Strategy); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateStrategy(s Strategy) error {
	switch s {
	case StrategyRoundRobin, StrategyRandom, StrategyLeastConnections:
		return nil
	default:
		return fmt.Errorf("stratégie inconnue %q, valeurs possibles: %q, %q, %q", s, StrategyRoundRobin, StrategyRandom, StrategyLeastConnections)
	}
}
