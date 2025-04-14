package metrics

import (
	"context"
	"fmt"
	"net/http"
	"ops_service/internal/config"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics represents Prometheus metrics
type Metrics struct {
	RequestsTotal      prometheus.Counter
	ResponseTime       prometheus.Histogram
	PickupPointsTotal  prometheus.Counter
	ReceiptsTotal      prometheus.Counter
	ProductsTotal      prometheus.Counter
}

// NewMetrics creates new Metrics
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		RequestsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "requests_total",
			Help:      "Total number of requests",
		}),
		ResponseTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "response_time_seconds",
			Help:      "Response time in seconds",
			Buckets:   prometheus.LinearBuckets(0.01, 0.01, 10), // 10ms to 100ms buckets
		}),
		PickupPointsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "pickup_points_created_total",
			Help:      "Total number of pickup points created",
		}),
		ReceiptsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "receipts_created_total",
			Help:      "Total number of receipts created",
		}),
		ProductsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "products_added_total",
			Help:      "Total number of products added",
		}),
	}
}

// ObserveResponseTime observes response time
func (m *Metrics) ObserveResponseTime(startTime time.Time) {
	m.ResponseTime.Observe(time.Since(startTime).Seconds())
}

// IncRequestsTotal increments total requests counter
func (m *Metrics) IncRequestsTotal() {
	m.RequestsTotal.Inc()
}

// IncPickupPointsTotal increments total pickup points counter
func (m *Metrics) IncPickupPointsTotal() {
	m.PickupPointsTotal.Inc()
}

// IncReceiptsTotal increments total receipts counter
func (m *Metrics) IncReceiptsTotal() {
	m.ReceiptsTotal.Inc()
}

// IncProductsTotal increments total products counter
func (m *Metrics) IncProductsTotal() {
	m.ProductsTotal.Inc()
}

// PrometheusHandler handles Prometheus metrics
type PrometheusHandler struct {
	metrics *Metrics
	server  *http.Server
}

// NewPrometheusHandler creates a new PrometheusHandler
func NewPrometheusHandler(metrics *Metrics, cfg config.PrometheusConfig) *PrometheusHandler {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}
	
	return &PrometheusHandler{
		metrics: metrics,
		server:  server,
	}
}

// Run runs the Prometheus server
func (h *PrometheusHandler) Run() error {
	return h.server.ListenAndServe()
}

// Shutdown gracefully shuts down the Prometheus server
func (h *PrometheusHandler) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
} 