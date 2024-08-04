package metric

import (
	"kleos/metrics/instrumentation"
	"strconv"
)

type Publisher interface {
	PublishMetrics(metricAttributes map[string]interface{}, metricType instrumentation.MetricType)
}

type PublisherImpl struct {
}

func NewMetricPublisher() *PublisherImpl {
	return &PublisherImpl{}
}

func (amp *PublisherImpl) PublishMetrics(metricAttributes instrumentation.MetricAttributes, metricType instrumentation.MetricType) {
	switch metricType {
	case instrumentation.ApiMetrics:
		publishApiMetric(metricAttributes.ApiMetric)
		return
	case instrumentation.ClientHttpCallMetrics:
		publishClientHttpCallMetric(metricAttributes.ClientHttpCallMetric)
		return
	default:
		return
	}
}

func publishApiMetric(apiMetrics instrumentation.ApiMetric) {
	status := strconv.Itoa(apiMetrics.ResponseCode)
	duration := float64(apiMetrics.DurationInMs)
	ApiRequestCounter.WithLabelValues(apiMetrics.Url, status).Inc()
	ApiRequestLatencyHistogram.WithLabelValues(apiMetrics.Url, status).Observe(duration)
}

func publishClientHttpCallMetric(clientHttpCallMetric instrumentation.ClientHttpCallMetric) {
	status := strconv.Itoa(clientHttpCallMetric.ResponseCode)
	duration := float64(clientHttpCallMetric.DurationInMs)
	HttpCallRequestCounter.WithLabelValues(clientHttpCallMetric.Url, status).Inc()
	HttpCallRequestLatencyHistogram.WithLabelValues(clientHttpCallMetric.Url, status).Observe(duration)
}
