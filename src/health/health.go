package health

import (
	"fmt"
	"liebe/src/utils"
	"log"
	"net/http"
	"time"
)

// Periodically runs health checks on all potential upstreams and updates the healthy list.
// - potentialUpstreams: full list of configured upstreams
// - updateHealthy: callback to update the list of healthy upstreams
func StartHealthCheck(potentialUpstreams []string, updateHealthy func([]string), path string, interval, timeout time.Duration) {
	// Run an initial pass immediately.
	go runHealthCheck(potentialUpstreams, updateHealthy, path, timeout)

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			runHealthCheck(potentialUpstreams, updateHealthy, path, timeout)
		}
	}()
}

// Performs a full health-check cycle and calls updateHealthy with the new healthy list.
func runHealthCheck(potentialUpstreams []string, updateHealthy func([]string), path string, timeout time.Duration) {
	client := &http.Client{Timeout: timeout}

	newHealthy := make([]string, 0, len(potentialUpstreams))
	unhealthyLogs := make([]string, 0)

	for _, upstream := range potentialUpstreams {
		ok, cause := checkHealth(client, upstream, path)
		if ok {
			newHealthy = append(newHealthy, upstream)
		} else if cause != "" {
			unhealthyLogs = append(unhealthyLogs, fmt.Sprintf("- %s -> cause: %s", upstream, cause))
		}
	}

	updateHealthy(newHealthy)

	if len(unhealthyLogs) > 0 {
		log.Printf("%d unhealthy upstream(s):\n%s", len(unhealthyLogs), utils.JoinLines(unhealthyLogs))
	}
}

// Performs a simple health check on a single upstream.
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

