package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ApiRequestCounter = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kleos_api_request_total",
		Help: "api request counter",
	},
	[]string{"url", "response_code"},
)

var ApiRequestLatencyHistogram = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "kleos_api_request_latency_histogram",
		Help: "api latency histogram",
	},
	[]string{"url", "response_code"},
)

var HttpCallRequestCounter = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kleos_http_call_request_total",
		Help: "http call request counter",
	},
	[]string{"url", "response_code"},
)

var HttpCallRequestLatencyHistogram = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "kleos_http_call_request_latency_histogram",
		Help: "http call latency histogram",
	},
	[]string{"url", "response_code"},
)
