package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(DeviceConnectionCount)
	prometheus.MustRegister(SSHConnectionCount)
	prometheus.MustRegister(SSHConnectionAttempts)
	prometheus.MustRegister(DeadAPIObjectCount)
	prometheus.MustRegister(DeadAPIQueries)
}

var (
	DeviceConnectionCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "device_connetion_count",
		Help: "Number of device connection being served",
	})

	DeadAPIObjectCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dead_api_object_count",
		Help: "Number of dead api object count",
	},
		[]string{"type"})

	DeadAPIQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dead_api_object_req_total",
			Help: "Number of dead object requests total",
		},
		[]string{"result", "type"})

	SSHConnectionCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "device_ssh_connetion_count",
			Help: "Number of current open ssh connections being served",
		},
		[]string{"source"})

	SSHConnectionAttempts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "device_ssh_connetion_attempts",
			Help: "Number of ssh connections attempts",
		},
		[]string{"result", "type"})
)
