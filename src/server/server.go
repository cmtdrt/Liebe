package server

import (
	"liebe/src/config"
	"liebe/src/strategy"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	potentialUpstreams []string
	mu                 sync.RWMutex
	healthyUpstreams   []string
	strategy           strategy.StrategyChooser
}

func NewLoadBalancer(cfg *config.Config) *LoadBalancer {
	return &LoadBalancer{
		potentialUpstreams: cfg.Upstreams,
		healthyUpstreams:   nil,
		strategy:           strategy.NewStrategyChooser(cfg.Strategy),
	}
}

// Returns the list of configured upstreams (healthy or not)
func (lb *LoadBalancer) PotentialUpstreams() []string {
	return lb.potentialUpstreams
}

// Updates the list of healthy upstreams.
func (lb *LoadBalancer) UpdateHealthy(healthy []string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.healthyUpstreams = healthy
}

// Selects the next upstream using the configured strategy.
func (lb *LoadBalancer) nextUpstream() (string, error) {
	lb.mu.RLock()
	healthyCopy := make([]string, len(lb.healthyUpstreams))
	copy(healthyCopy, lb.healthyUpstreams)
	lb.mu.RUnlock()

	return lb.strategy.Next(healthyCopy)
}

// Proxies incoming requests to a healthy upstream.
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target, err := lb.nextUpstream()
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
