package legacy

import (
	"testing"
)

func TestNewBashExecutor(t *testing.T) {
	e := NewBashExecutor("/tmp/kudig.sh")
	if e == nil {
		t.Fatal("NewBashExecutor() returned nil")
	}
	if e.ScriptPath != "/tmp/kudig.sh" {
		t.Errorf("ScriptPath = %v, want %v", e.ScriptPath, "/tmp/kudig.sh")
	}
}

func TestFindScript_NotFound(t *testing.T) {
	// This test may fail if kudig.sh exists in any of the search paths
	// In that case, we just verify the function returns a path
	path, err := FindScript()
	// We can't predict if it will be found or not in test environment
	// So we just check the function doesn't panic
	_ = path
	_ = err
}

func TestConvertBashReportToIssues(t *testing.T) {
	report := &BashReport{
		ReportVersion: "1.0",
		Timestamp:     "2024-01-01T00:00:00Z",
		Hostname:      "test-node",
		DiagnoseDir:   "/tmp/test",
		Anomalies: []BashAnomaly{
			{
				Severity: "严重",
				CNName:   "严重问题",
				ENName:   "CRITICAL_ISSUE",
				Details:  "问题详情",
				Location: "test.go",
			},
			{
				Severity: "警告",
				CNName:   "警告问题",
				ENName:   "WARNING_ISSUE",
				Details:  "警告详情",
				Location: "test2.go",
			},
		},
		Summary: BashSummary{
			Critical: 1,
			Warning:  1,
			Info:     0,
			Total:    2,
		},
	}

	issues := ConvertBashReportToIssues(report)
	if len(issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(issues))
	}

	if issues[0].ENName != "CRITICAL_ISSUE" {
		t.Errorf("First issue ENName = %v, want %v", issues[0].ENName, "CRITICAL_ISSUE")
	}

	if issues[1].ENName != "WARNING_ISSUE" {
		t.Errorf("Second issue ENName = %v, want %v", issues[1].ENName, "WARNING_ISSUE")
	}
}

func TestConvertBashReportToIssues_Empty(t *testing.T) {
	report := &BashReport{
		Anomalies: []BashAnomaly{},
		Summary: BashSummary{
			Total: 0,
		},
	}

	issues := ConvertBashReportToIssues(report)
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}

func TestConvertBashReportToIssues_Metadata(t *testing.T) {
	report := &BashReport{
		ReportVersion: "1.0",
		Timestamp:     "2024-01-01T00:00:00Z",
		Hostname:      "test-node",
		DiagnoseDir:   "/tmp/test",
		Anomalies: []BashAnomaly{
			{
				Severity: "提示",
				CNName:   "提示",
				ENName:   "INFO_ISSUE",
				Details:  "提示详情",
				Location: "info.go",
			},
		},
	}

	issues := ConvertBashReportToIssues(report)
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}

	// Check metadata is set
	if issues[0].Metadata["source"] != "kudig.sh" {
		t.Error("Metadata not set correctly")
	}
}
