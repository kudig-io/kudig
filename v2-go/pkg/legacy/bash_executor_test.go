package legacy

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
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
	path, err := FindScript()
	_ = path
	_ = err
}

func TestFindScript_FoundInDir(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "kudig.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/bash\necho ok"), 0755); err != nil {
		t.Fatal(err)
	}
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	path, err := FindScript()
	if err != nil {
		t.Fatalf("FindScript() error = %v", err)
	}
	if path == "" {
		t.Error("FindScript() should return a path")
	}
}

func TestConvertBashReportToIssues(t *testing.T) {
	report := &BashReport{
		ReportVersion: "1.0",
		Timestamp:     "2024-01-01T00:00:00Z",
		Hostname:      "test-node",
		DiagnoseDir:   "/tmp/test",
		Anomalies: []BashAnomaly{
			{Severity: "严重", CNName: "严重问题", ENName: "CRITICAL_ISSUE", Details: "问题详情", Location: "test.go"},
			{Severity: "警告", CNName: "警告问题", ENName: "WARNING_ISSUE", Details: "警告详情", Location: "test2.go"},
		},
		Summary: BashSummary{Critical: 1, Warning: 1, Total: 2},
	}
	issues := ConvertBashReportToIssues(report)
	if len(issues) != 2 {
		t.Fatalf("Expected 2 issues, got %d", len(issues))
	}
	if issues[0].ENName != "CRITICAL_ISSUE" {
		t.Errorf("First issue ENName = %v, want CRITICAL_ISSUE", issues[0].ENName)
	}
	if issues[0].AnalyzerName != "legacy.bash" {
		t.Errorf("AnalyzerName = %v, want legacy.bash", issues[0].AnalyzerName)
	}
	if issues[1].ENName != "WARNING_ISSUE" {
		t.Errorf("Second issue ENName = %v, want WARNING_ISSUE", issues[1].ENName)
	}
}

func TestConvertBashReportToIssues_Empty(t *testing.T) {
	report := &BashReport{Anomalies: []BashAnomaly{}}
	issues := ConvertBashReportToIssues(report)
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}

func TestConvertBashReportToIssues_Metadata(t *testing.T) {
	report := &BashReport{
		Anomalies: []BashAnomaly{
			{Severity: "提示", CNName: "提示", ENName: "INFO_ISSUE", Details: "提示详情"},
		},
	}
	issues := ConvertBashReportToIssues(report)
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].Metadata["source"] != "kudig.sh" {
		t.Error("Metadata source not set correctly")
	}
}

func TestExecute_ScriptNotFound(t *testing.T) {
	e := NewBashExecutor("/nonexistent/kudig.sh")
	issues, report, err := e.Execute(context.Background(), "/tmp/diag", false)
	if err == nil {
		t.Error("Expected error for nonexistent script")
	}
	if issues != nil {
		t.Error("Expected nil issues on error")
	}
	if report != nil {
		t.Error("Expected nil report on error")
	}
}

func TestExecute_NoOutput(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "kudig.sh")
	os.WriteFile(scriptPath, []byte("#!/bin/bash\nexit 0"), 0755)

	e := NewBashExecutor(scriptPath)
	issues, report, err := e.Execute(context.Background(), "/tmp/diag", false)
	if err == nil {
		t.Error("Expected error for no output")
	}
	if issues != nil || report != nil {
		t.Error("Expected nil results for no output")
	}
}

func TestExecute_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "kudig.sh")
	os.WriteFile(scriptPath, []byte("#!/bin/bash\necho 'not json'"), 0755)

	e := NewBashExecutor(scriptPath)
	_, _, err := e.Execute(context.Background(), "/tmp/diag", false)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestExecute_ValidOutput(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "kudig.sh")
	output := `{"report_version":"1.0","timestamp":"2024-01-01T00:00:00Z","hostname":"test","diagnose_dir":"/tmp","anomalies":[{"severity":"严重","cn_name":"测试","en_name":"TEST","details":"test details","location":""}],"summary":{"critical":1,"warning":0,"info":0,"total":1}}`
	os.WriteFile(scriptPath, []byte("#!/bin/bash\necho '"+output+"'"), 0755)

	e := NewBashExecutor(scriptPath)
	issues, report, err := e.Execute(context.Background(), "/tmp/diag", false)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].ENName != "TEST" {
		t.Errorf("ENName = %v, want TEST", issues[0].ENName)
	}
	if report == nil {
		t.Fatal("Expected non-nil report")
	}
	if report.Summary.Total != 1 {
		t.Errorf("Summary.Total = %d, want 1", report.Summary.Total)
	}
}

func TestExecute_Verbose(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "kudig.sh")
	output := `{"report_version":"1.0","timestamp":"","hostname":"","diagnose_dir":"","anomalies":[],"summary":{"total":0}}`
	os.WriteFile(scriptPath, []byte("#!/bin/bash\necho '"+output+"'"), 0755)

	e := NewBashExecutor(scriptPath)
	_, _, err := e.Execute(context.Background(), "/tmp/diag", true)
	if err != nil {
		t.Fatalf("Execute() with verbose error = %v", err)
	}
}

func TestExecute_CancelledContext(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "kudig.sh")
	os.WriteFile(scriptPath, []byte("#!/bin/bash\nsleep 60"), 0755)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	e := NewBashExecutor(scriptPath)
	_, _, err := e.Execute(ctx, "/tmp/diag", false)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
}

func TestNewLegacyCollector_WithScriptPath(t *testing.T) {
	c, err := NewLegacyCollector("/tmp/kudig.sh")
	if err != nil {
		t.Fatalf("NewLegacyCollector() error = %v", err)
	}
	if c == nil {
		t.Fatal("NewLegacyCollector() returned nil")
	}
	if c.executor.ScriptPath != "/tmp/kudig.sh" {
		t.Errorf("ScriptPath = %v, want /tmp/kudig.sh", c.executor.ScriptPath)
	}
}

func TestNewLegacyCollector_EmptyPath_NotFound(t *testing.T) {
	_, err := NewLegacyCollector("")
	if err == nil {
		t.Error("Expected error when kudig.sh not found")
	}
}

func TestLegacyCollector_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "kudig.sh")
	output := `{"report_version":"1.0","timestamp":"","hostname":"","diagnose_dir":"","anomalies":[{"severity":"警告","cn_name":"测试","en_name":"WARN","details":"detail"}],"summary":{"warning":1,"total":1}}`
	os.WriteFile(scriptPath, []byte("#!/bin/bash\necho '"+output+"'"), 0755)

	c, _ := NewLegacyCollector(scriptPath)
	issues, err := c.Execute(context.Background(), "/tmp/diag", false)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
}

func TestLegacyCollector_GetReport(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "kudig.sh")
	output := `{"report_version":"1.0","timestamp":"2024-01-01","hostname":"h","diagnose_dir":"/d","anomalies":[],"summary":{"total":0}}`
	os.WriteFile(scriptPath, []byte("#!/bin/bash\necho '"+output+"'"), 0755)

	c, _ := NewLegacyCollector(scriptPath)
	report, err := c.GetReport(context.Background(), "/tmp/diag", false)
	if err != nil {
		t.Fatalf("GetReport() error = %v", err)
	}
	if report == nil {
		t.Fatal("Expected non-nil report")
	}
	if report.ReportVersion != "1.0" {
		t.Errorf("ReportVersion = %v, want 1.0", report.ReportVersion)
	}
}

func TestBashReport_JSONRoundTrip(t *testing.T) {
	report := BashReport{
		ReportVersion: "2.0",
		Timestamp:     "2024-06-01",
		Hostname:      "roundtrip-node",
		DiagnoseDir:   "/tmp/rt",
		Anomalies: []BashAnomaly{
			{Severity: "严重", CNName: "测试", ENName: "RT_TEST", Details: "detail", Location: "loc"},
		},
		Summary: BashSummary{Critical: 1, Total: 1},
	}
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	var decoded BashReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if decoded.ReportVersion != "2.0" {
		t.Errorf("ReportVersion = %v, want 2.0", decoded.ReportVersion)
	}
	if len(decoded.Anomalies) != 1 {
		t.Errorf("Anomalies length = %d, want 1", len(decoded.Anomalies))
	}
}
