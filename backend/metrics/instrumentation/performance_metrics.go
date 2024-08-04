package instrumentation

type MetricType string

const (
	ApiMetrics            MetricType = "API_METRICS"
	ClientHttpCallMetrics MetricType = "CLIENT_HTTP_CALL_METRICS"
)

type ApiMetric struct {
	Url           string `json:"url,omitempty"`
	Method        string `json:"method,omitempty"`
	ResponseCode  int    `json:"response_code,omitempty"`
	BytesSent     int    `json:"bytes_sent,omitempty"`
	BytesReceived int64  `json:"bytes_received,omitempty"`
	StartTime     int64  `json:"start_time,omitempty"`
	EndTime       int64  `json:"end_time,omitempty"`
	DurationInMs  int64  `json:"duration_in_ms,omitempty"`
	ErrorType     string `json:"error_type,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

type ClientHttpCallMetric struct {
	Url          string `json:"url,omitempty"`
	ResponseCode int    `json:"response_code,omitempty"`
	StartTime    int64  `json:"start_time,omitempty"`
	EndTime      int64  `json:"end_time,omitempty"`
	DurationInMs int64  `json:"duration_in_ms,omitempty"`
}

type MetricAttributes struct {
	ApiMetric            ApiMetric
	ClientHttpCallMetric ClientHttpCallMetric
}
