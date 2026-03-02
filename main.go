package main

import (
	"fmt"
	"liebe/src/config"
	"liebe/src/utils"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type StrategyChooser interface {
	Next(healthy []string) (string, error)
}

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

type randomStrategy struct {
	rnd *rand.Rand
}

func (s *randomStrategy) Next(healthy []string) (string, error) {
	if len(healthy) == 0 {
		return "", http.ErrServerClosed
	}
	return healthy[s.rnd.Intn(len(healthy))], nil
}

type LoadBalancer struct {
	potentialUpstreams []string

	mu               sync.RWMutex
	healthyUpstreams []string

	strategy StrategyChooser
}

func NewLoadBalancer(cfg *config.Config) *LoadBalancer {
	var strat StrategyChooser
	switch cfg.Strategy {
	case config.StrategyRoundRobin:
		strat = &roundRobinStrategy{}
	case config.StrategyRandom:
		strat = &randomStrategy{rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
	default:
		panic(fmt.Sprintf("Unknown strategy %s", cfg.Strategy)) // Should never happen
	}

	return &LoadBalancer{
		potentialUpstreams: cfg.Upstreams,
		healthyUpstreams:   nil,
		strategy:           strat,
	}
}

// lance un health-check de tous les potentialUpstreams toutes les "interval" secondes
func (lb *LoadBalancer) StartHealthCheck(path string, interval, timeout time.Duration) {
	// première passe immédiate
	go lb.runHealthCheck(path, timeout)

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			lb.runHealthCheck(path, timeout)
		}
	}()
}

// health-check complet : met à jour la liste des healthyUpstreams
func (lb *LoadBalancer) runHealthCheck(path string, timeout time.Duration) {
	client := &http.Client{Timeout: timeout}

	newHealthy := make([]string, 0, len(lb.potentialUpstreams))
	unhealthyLogs := make([]string, 0)
	for _, upstream := range lb.potentialUpstreams {
		ok, cause := checkHealth(client, upstream, path)
		if ok {
			newHealthy = append(newHealthy, upstream)
		} else if cause != "" {
			unhealthyLogs = append(unhealthyLogs, fmt.Sprintf("- %s -> cause: %s", upstream, cause))
		}
	}

	lb.mu.Lock()
	lb.healthyUpstreams = newHealthy
	lb.mu.Unlock()

	if len(unhealthyLogs) > 0 {
		log.Printf("%d unhealthy upstream(s):\n%s", len(unhealthyLogs), utils.JoinLines(unhealthyLogs))
	}
}

// checks the health of an upstream
func checkHealth(client *http.Client, upstream, path string) (bool, string) {
	resp, err := client.Get(upstream + path)
	if err != nil {
		return false, err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Sprintf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return true, ""
}

func (lb *LoadBalancer) selectNextUpstream() (string, error) {
	lb.mu.RLock()
	healthyCopy := make([]string, len(lb.healthyUpstreams))
	copy(healthyCopy, lb.healthyUpstreams)
	lb.mu.RUnlock()

	return lb.strategy.Next(healthyCopy)
}

// proxy vers le bon upstream sain
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target, err := lb.selectNextUpstream()
	if err != nil || target == "" {
		http.Error(w, "aucun upstream sain disponible", http.StatusServiceUnavailable)
		return
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "upstream invalide", http.StatusInternalServerError)
		return
	}

	log.Printf("Incoming %s request on %s routed to upstream %s", r.Method, r.URL.Path, targetURL.String())

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Println("proxy error:", err)
		http.Error(w, "erreur en appelant l'upstream", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)
}

func main() {
	cfg, err := config.LoadConfig("liebe-config.json")
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}
	log.Println("Config loaded successfully:", cfg)

	lb := NewLoadBalancer(cfg)

	interval := time.Duration(cfg.HealthCheck.Interval) * time.Second
	timeout := time.Duration(cfg.HealthCheck.Timeout) * time.Second
	lb.StartHealthCheck(cfg.HealthCheck.Path, interval, timeout)

	addr := ":8080"
	log.Println("Load balancer listening on", addr)
	if err := http.ListenAndServe(addr, lb); err != nil {
		log.Fatal(err)
	}
}
