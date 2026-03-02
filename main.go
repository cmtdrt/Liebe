package main

import (
	"liebe/src/config"
	"liebe/src/health"
	"liebe/src/server"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("liebe-config.json")
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}
	log.Println("Config loaded successfully:", cfg)

	lb := server.NewLoadBalancer(cfg)

	interval := time.Duration(cfg.HealthCheck.Interval) * time.Second
	timeout := time.Duration(cfg.HealthCheck.Timeout) * time.Second

	health.StartHealthCheck(
		lb.PotentialUpstreams(),
		lb.UpdateHealthy,
		cfg.HealthCheck.Path,
		interval,
		timeout,
	)

	addr := ":8080"
	log.Println("Load balancer listening on", addr)
	if err := http.ListenAndServe(addr, lb); err != nil {
		log.Fatal(err)
	}
}

