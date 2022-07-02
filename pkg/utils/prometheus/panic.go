package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(GOPanicCounter)
}

var (
	// TotalRequests is the total number of HTTP requests made.
	GOPanicCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_panic_counter",
			Help: "Number of panics in go code.",
		},
		[]string{"message"},
	)
)
