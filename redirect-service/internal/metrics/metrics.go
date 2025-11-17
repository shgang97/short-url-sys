package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RedisHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redirect_service_redis_hits_total",
			Help: "Total number of Redis hits",
		},
	)

	RedismissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redirect_service_redis_misses_total",
			Help: "Total number of Redis misses",
		},
	)

	RedisErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirect_service_redis_errors_total",
			Help: "Total number of Redis errors",
		},
		[]string{"operation"},
	)

	CacheHitRatio = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redirect_service_redis_hit_ratio",
			Help: "Cache hit ratio",
		},
	)

	RedisHealthStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redirect_service_redis_health_status",
			Help: "Redis health status (1 = healthy, 0 = unhealthy)",
		},
	)
)
