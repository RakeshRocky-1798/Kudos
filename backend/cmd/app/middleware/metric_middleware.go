package middleware

import (
	"github.com/gin-gonic/gin"
	"kleos/metrics/instrumentation"
	"kleos/metrics/metric"
	"time"
)

func MetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		duration := endTime.Sub(startTime)

		metricsPublisher := metric.NewMetricPublisher()
		apiMetrics := instrumentation.ApiMetric{
			Url:          c.Request.URL.Path,
			Method:       c.Request.Method,
			ResponseCode: c.Writer.Status(),
			StartTime:    startTime.Unix(),
			EndTime:      endTime.Unix(),
			DurationInMs: duration.Milliseconds(),
			BytesSent:    c.Writer.Size(),
		}

		metricsPublisher.PublishMetrics(instrumentation.MetricAttributes{ApiMetric: apiMetrics},
			instrumentation.ApiMetrics)
	}
}
