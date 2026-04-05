// Package metrics provides Prometheus metrics for kudig diagnostic operations.
package metrics

import (
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kudig/kudig/pkg/types"
)

var (
	// namespace for all kudig metrics
	ns = "kudig"
)

// Registry holds all kudig metrics.
var Registry = prometheus.NewRegistry()

var (
	// DiagnosesTotal counts the total number of diagnostic runs.
	DiagnosesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "diagnoses_total",
			Help:      "Total number of diagnostic runs",
		},
		[]string{"mode", "status"},
	)

	// DiagnosisDuration records the duration of diagnostic runs.
	DiagnosisDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ns,
			Name:      "diagnosis_duration_seconds",
			Help:      "Duration of diagnostic runs in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"mode"},
	)

	// IssuesTotal records the number of issues found per severity and category.
	IssuesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "issues_total",
			Help:      "Number of issues found by severity and category",
		},
		[]string{"severity", "category"},
	)

	// AnalyzersExecutedTotal counts the number of analyzers executed.
	AnalyzersExecutedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ns,
			Name:      "analyzers_executed_total",
			Help:      "Total number of analyzers executed",
		},
		[]string{"mode"},
	)
)

func init() {
	Registry.MustRegister(DiagnosesTotal)
	Registry.MustRegister(DiagnosisDuration)
	Registry.MustRegister(IssuesTotal)
	Registry.MustRegister(AnalyzersExecutedTotal)
}

// Server wraps an HTTP server that exposes Prometheus metrics.
type Server struct {
	addr string
}

// NewServer creates a new metrics server on the given address.
func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

// Start begins serving metrics on the configured address.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	//nolint:gosec // metrics server intentionally does not use timeouts in this minimal implementation
	return http.ListenAndServe(s.addr, mux)
}

// Addr returns the server address.
func (s *Server) Addr() string {
	return s.addr
}

// RecordDiagnosis records a completed diagnostic run.
func RecordDiagnosis(mode types.DataMode, status string, duration time.Duration) {
	modeStr := mode.String()
	DiagnosesTotal.WithLabelValues(modeStr, status).Inc()
	DiagnosisDuration.WithLabelValues(modeStr).Observe(duration.Seconds())
}

// RecordIssues updates the issue gauges based on the provided issues.
func RecordIssues(issues []types.Issue) {
	// Reset gauges so stale values are cleared
	IssuesTotal.Reset()

	// Count issues by severity and category
	counts := make(map[string]map[string]int)
	for _, issue := range issues {
		sev := severityString(issue.Severity)
		cat := issueCategory(issue)
		if counts[sev] == nil {
			counts[sev] = make(map[string]int)
		}
		counts[sev][cat]++
	}

	for sev, catMap := range counts {
		for cat, count := range catMap {
			IssuesTotal.WithLabelValues(sev, cat).Set(float64(count))
		}
	}
}

// RecordAnalyzers records the number of analyzers executed.
func RecordAnalyzers(mode types.DataMode, count int) {
	AnalyzersExecutedTotal.WithLabelValues(mode.String()).Add(float64(count))
}

func severityString(s types.Severity) string {
	switch s {
	case types.SeverityCritical:
		return "critical"
	case types.SeverityWarning:
		return "warning"
	case types.SeverityInfo:
		return "info"
	default:
		return "unknown"
	}
}

func issueCategory(issue types.Issue) string {
	if issue.AnalyzerName != "" {
		// Extract category from analyzer name (e.g., "system.cpu" -> "system")
		parts := strings.Split(issue.AnalyzerName, ".")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return "unknown"
}
