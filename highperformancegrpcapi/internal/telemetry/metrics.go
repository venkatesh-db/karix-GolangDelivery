package telemetry

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	ActiveSessions     prometheus.Gauge
	IngressMessages    *prometheus.CounterVec
	EgressMessages     *prometheus.CounterVec
	DroppedMessages    *prometheus.CounterVec
	SessionDuration    prometheus.Summary
	QueueDepth         prometheus.Summary
	HeartbeatMissCount prometheus.Counter
}

func New(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		ActiveSessions: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "ride_stream_active_sessions",
			Help: "Current live sessions across transports",
		}),
		IngressMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "ride_stream_ingress_total",
			Help: "Count of inbound envelopes by type",
		}, []string{"type", "transport"}),
		EgressMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "ride_stream_egress_total",
			Help: "Count of outbound envelopes by type",
		}, []string{"type", "transport"}),
		DroppedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "ride_stream_dropped_total",
			Help: "Number of envelopes dropped due to backpressure",
		}, []string{"reason"}),
		SessionDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Name:       "ride_stream_session_duration_seconds",
			Help:       "Session length distribution",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		QueueDepth: prometheus.NewSummary(prometheus.SummaryOpts{
			Name:       "ride_stream_queue_depth",
			Help:       "Snapshot of outbound queue depth when sampled",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		HeartbeatMissCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ride_stream_heartbeat_miss_total",
			Help: "Number of heartbeat windows missed",
		}),
	}

	reg.MustRegister(
		m.ActiveSessions,
		m.IngressMessages,
		m.EgressMessages,
		m.DroppedMessages,
		m.SessionDuration,
		m.QueueDepth,
		m.HeartbeatMissCount,
	)

	return m
}

func Handler(reg prometheus.Gatherer) http.Handler {
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{Timeout: 5 * time.Second})
}
