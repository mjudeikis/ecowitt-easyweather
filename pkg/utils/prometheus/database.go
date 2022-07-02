package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(DatabaseRequestCount)
}

var (
	// DatabaseCalls is the total number of database calls.
	DatabaseRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_request_count",
			Help: "Number of database request calls. With cache misses",
		},
		[]string{"method", "cached"},
	)
)
