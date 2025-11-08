package circuitbreaker

import (
	"log"
	"time"

	"github.com/sony/gobreaker"
)

type RedisCircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

type RedisCircuitBreakerConfig struct {
	Name        string        `mapstructure:"name"`
	MaxRequests uint32        `mapstructure:"max_requests"` // 熔断器打开前的最大请求数
	Timeout     time.Duration `mapstructure:"timeout"`      // 熔断器打开前的超时时间
	Interval    time.Duration `mapstructure:"interval"`     // 熔断器打开后，进入半开状态前的等待时间
	Threshold   float64       `mapstructure:"threshold"`    // 熔断器打开的条件：失败率 >= Threshold
}

func NewRedisCircuitBreaker(cfg *RedisCircuitBreakerConfig) *RedisCircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Timeout:     cfg.Timeout * time.Second,
		Interval:    cfg.Interval * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			total := counts.Requests
			if total < 5 {
				return false
			}
			failureRatio := float64(counts.TotalFailures) / float64(total)
			return failureRatio >= (cfg.Threshold / 100)
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker %s state change: %s -> %s", name, from, to)
		},
	}
	return &RedisCircuitBreaker{
		cb: gobreaker.NewCircuitBreaker(settings),
	}
}

// Execute 执行Redis操作，如果熔断器打开则返回错误
func (rcb *RedisCircuitBreaker) Execute(operation func() (interface{}, error)) (interface{}, error) {
	req := func() (interface{}, error) {
		return operation()
	}
	return rcb.cb.Execute(req)
}
