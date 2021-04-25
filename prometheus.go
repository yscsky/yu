package yu

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// prmethous监控变量
var (
	RequestCounter   *prometheus.CounterVec
	LatencyHistogram *prometheus.HistogramVec
	DurationsSummary *prometheus.SummaryVec
)

// InitPrometheus 初始化Prometheus
func InitPrometheus(name string) {
	if strings.Contains(name, "-") {
		name = strings.ReplaceAll(name, "-", "_")
	}
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "SHORT",
			Subsystem: name,
			Name:      "RequestCounter",
			Help:      "Requests Count",
		},
		[]string{"service"},
	)

	LatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "SHORT",
			Subsystem: name,
			Name:      "Latency",
			Help:      "Requests Latency Histogram",
			Buckets:   []float64{0, 1e+5 /*100 us*/, 1e+6 /*1 ms*/, 1e+7 /*10 ms*/, 1e+8 /*100 ms*/, 1e+9 /*1 s*/, 1e+10 /*10 s*/},
		},
		[]string{"service"},
	)

	DurationsSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "SHORT",
			Subsystem:  name,
			Name:       "LatencySummary",
			Help:       "Requests Latency Summary",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service"},
	)
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(LatencyHistogram)
	prometheus.MustRegister(DurationsSummary)
}

// PromMetrics prometheus metrics 中间件
func PromMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(start time.Time) {
			d := float64(time.Since(start).Nanoseconds())
			LatencyHistogram.WithLabelValues(c.Request.URL.Path).Observe(d)
			DurationsSummary.WithLabelValues(c.Request.URL.Path).Observe(d)
			RequestCounter.WithLabelValues(c.Request.URL.Path).Add(1)
		}(time.Now())
		c.Next()
	}
}

// PromethousHandler 启动promethous http监听
func PromethousHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}

// PrometheusInterceptor promethous grpc中间件
func PrometheusInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func(start time.Time) {
		d := float64(time.Since(start).Nanoseconds())
		LatencyHistogram.WithLabelValues(info.FullMethod).Observe(d)
		DurationsSummary.WithLabelValues(info.FullMethod).Observe(d)
		RequestCounter.WithLabelValues(info.FullMethod).Add(1)
	}(time.Now())
	return handler(ctx, req)
}
