package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(BytesReceivedCounter)
	prometheus.MustRegister(BytesTransferredCounter)
}

var (
	// TotalRequests is the total number of HTTP requests made.
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path", "status", "method", "user_agent"},
	)

	// HttpDuration is a Prometheus Histogram metric for the HTTP request latency.
	HttpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path", "method"})

	BytesTransferredCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "controller_bytes_transferred_total",
			Help: "Total number of bytes transferred.",
		},
		[]string{"path"},
	)

	BytesReceivedCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "controller_bytes_received_total",
			Help: "Total number of bytes received.",
		},
		[]string{"path"},
	)
)
