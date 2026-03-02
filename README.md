### Liebe

Liebe is a lightweight, configurable load balancer that distributes traffic across multiple identical API instances, with built‑in health checks and simple routing strategies.

### Configuration

Configuration lives in `liebe-config.json` at the project root:

```json
{
  "health_check": {
    "path": "/health",
    "interval": "5",
    "timeout": "2"
  },
  "strategy": "round_robin",
  "upstreams": [
    "http://localhost:8081",
    "http://localhost:8082",
    "http://localhost:8083"
  ]
}
```

- **health_check.path**: endpoint called on each upstream (e.g. `/health`).
- **health_check.interval**: every X seconds, Liebe checks all upstreams.
- **health_check.timeout**: maximum time allowed for a response to be considered valid.
- **strategy**: load‑balancing strategy (`"round_robin"` or `"random"`).
- **upstreams**: list of API instances that will receive traffic.

### How it works

1. On startup, Liebe loads `liebe-config.json` and validates the strategy.
2. In a loop, a health check calls the configured endpoint on every upstream:
   - if it returns `200` before `timeout`, the upstream is marked **healthy**;
   - otherwise it is marked **unhealthy** and detailed in the logs.
3. For each incoming request on `:8080`, Liebe picks a **healthy** upstream according to the strategy and proxies the request to it.