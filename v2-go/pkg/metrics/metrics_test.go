package metrics

import (
	"testing"
	"time"

	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/kudig/kudig/pkg/types"
)

func TestRecordDiagnosis(t *testing.T) {
	// Clear previous test data
	DiagnosesTotal.Reset()
	DiagnosisDuration.Reset()

	// Record a diagnosis
	RecordDiagnosis(types.ModeOnline, "success", 5*time.Second)

	// Check counter was incremented
	cnt, err := DiagnosesTotal.GetMetricWithLabelValues("online", "success")
	if err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	metric := &io_prometheus_client.Metric{}
	cnt.Write(metric)
	if metric.Counter.GetValue() != 1 {
		t.Errorf("Expected counter value 1, got %f", metric.Counter.GetValue())
	}
}

func TestRecordIssues(t *testing.T) {
	// Clear previous data
	IssuesTotal.Reset()

	// Create test issues
	issues := []types.Issue{
		{Severity: types.SeverityCritical, AnalyzerName: "system.cpu"},
		{Severity: types.SeverityWarning, AnalyzerName: "network.interface"},
		{Severity: types.SeverityCritical, AnalyzerName: "kubernetes.pod"},
		{Severity: types.SeverityInfo, AnalyzerName: "runtime.config"},
	}

	// Record issues
	RecordIssues(issues)

	// Verify metrics were recorded
	critical, err := IssuesTotal.GetMetricWithLabelValues("critical", "system")
	if err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	metric := &io_prometheus_client.Metric{}
	critical.Write(metric)
	if metric.Gauge.GetValue() != 1 {
		t.Errorf("Expected gauge value 1, got %f", metric.Gauge.GetValue())
	}

	// Check that we have the right number of unique category/severity combinations
	families, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	found := false
	for _, f := range families {
		if f.GetName() == "kudig_issues_total" {
			found = true
			if len(f.GetMetric()) < 3 {
				t.Errorf("Expected at least 3 metric combinations, got %d", len(f.GetMetric()))
			}
			break
		}
	}
	if !found {
		t.Error("kudig_issues_total metric should exist")
	}
}

func TestRecordAnalyzers(t *testing.T) {
	AnalyzersExecutedTotal.Reset()

	// Record analyzers
	RecordAnalyzers(types.ModeOnline, 5)
	RecordAnalyzers(types.ModeOnline, 3)

	// Verify
	cnt, err := AnalyzersExecutedTotal.GetMetricWithLabelValues("online")
	if err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	metric := &io_prometheus_client.Metric{}
	cnt.Write(metric)
	if metric.Counter.GetValue() != 8 {
		t.Errorf("Expected counter value 8, got %f", metric.Counter.GetValue())
	}
}

func TestSeverityString(t *testing.T) {
	tests := []struct {
		input    types.Severity
		expected string
	}{
		{types.SeverityCritical, "critical"},
		{types.SeverityWarning, "warning"},
		{types.SeverityInfo, "info"},
		{types.Severity(99), "unknown"},
	}

	for _, tc := range tests {
		result := severityString(tc.input)
		if result != tc.expected {
			t.Errorf("severityString(%v) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestIssueCategory(t *testing.T) {
	tests := []struct {
		analyzerName string
		expected     string
	}{
		{"system.cpu", "system"},
		{"network.interface", "network"},
		{"kubernetes.pod", "kubernetes"},
		{"single", "single"},
		{"", "unknown"},
	}

	for _, tc := range tests {
		issue := types.Issue{AnalyzerName: tc.analyzerName}
		result := issueCategory(issue)
		if result != tc.expected {
			t.Errorf("issueCategory(%q) = %s, expected %s", tc.analyzerName, result, tc.expected)
		}
	}
}
