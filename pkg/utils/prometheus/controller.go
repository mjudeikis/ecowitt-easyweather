package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(ControllerMetricsQueue)
}

var (
	ControllerMetricsQueue = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "controller_metrics_ingestion_queue",
		Help: "Number of metrics in the queue for ingestion",
	})
)
