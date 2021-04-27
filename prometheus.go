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
func InitPrometheus(space, name string) {
	if strings.Contains(name, "-") {
		name = strings.ReplaceAll(name, "-", "_")
	}
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: space,
			Subsystem: name,
			Name:      "RequestCounter",
			Help:      "Requests Count",
		},
		[]string{"service"},
	)

	LatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: space,
			Subsystem: name,
			Name:      "Latency",
			Help:      "Requests Latency Histogram",
			Buckets:   []float64{0, 1e+5 /*100 us*/, 1e+6 /*1 ms*/, 1e+7 /*10 ms*/, 1e+8 /*100 ms*/, 1e+9 /*1 s*/, 1e+10 /*10 s*/},
		},
		[]string{"service"},
	)

	DurationsSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  space,
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

// PromethousHandler 启动promethous http监听
func PromethousHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}

// PromeMetrics prometheus metrics 中间件
func PromeMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer PromeTrace(c.Request.URL.Path)()
		c.Next()
	}
}

// PromeUnaryInterceptor promethous grpc中间件
func PromeUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler) (resp interface{}, err error) {
			defer PromeTrace(info.FullMethod)()
			return handler(ctx, req)
		},
	)
}

// PromeTrace promethous数据记录
func PromeTrace(label string) func() {
	start := time.Now()
	return func() {
		d := float64(time.Since(start).Nanoseconds())
		LatencyHistogram.WithLabelValues(label).Observe(d)
		DurationsSummary.WithLabelValues(label).Observe(d)
		RequestCounter.WithLabelValues(label).Add(1)
	}
}
