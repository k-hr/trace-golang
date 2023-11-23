package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	StatusCodeLabel = "status_code"
	MethodLabel     = "method"
	RouteLabel      = "route"
)

var RequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "trace_app_prom_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	},
	[]string{StatusCodeLabel, MethodLabel, RouteLabel},
)

var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "trace_app_prom_duration",
		Help: "Duration of http requests.",
	},
	[]string{StatusCodeLabel, MethodLabel, RouteLabel},
)

func RegisterProm() {
	prometheus.MustRegister(
		RequestCount,
		RequestDuration,
	)
}
