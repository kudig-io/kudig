package process

import (
	"context"
	"testing"

	"github.com/kudig/kudig/pkg/types"
)

func TestNewKubeletAnalyzer(t *testing.T) {
	a := NewKubeletAnalyzer()
	if a == nil {
		t.Fatal("NewKubeletAnalyzer() returned nil")
	}
	if a.Name() != "process.kubelet" {
		t.Errorf("Name() = %v, want %v", a.Name(), "process.kubelet")
	}
}

func TestKubeletAnalyzerAnalyze_NoFile(t *testing.T) {
	a := NewKubeletAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}

func TestKubeletAnalyzerAnalyze_Running(t *testing.T) {
	a := NewKubeletAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"daemon_status/kubelet_status": []byte("Active: active (running)"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues for running service, got %d", len(issues))
	}
}

func TestKubeletAnalyzerAnalyze_Failed(t *testing.T) {
	a := NewKubeletAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"daemon_status/kubelet_status": []byte("Active: failed (Result: exit-code)"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].ENName != "KUBELET_SERVICE_DOWN" {
		t.Errorf("ENName = %v, want %v", issues[0].ENName, "KUBELET_SERVICE_DOWN")
	}
	if issues[0].Severity != types.SeverityCritical {
		t.Errorf("Severity = %v, want %v", issues[0].Severity, types.SeverityCritical)
	}
}

func TestKubeletAnalyzerAnalyze_Stopped(t *testing.T) {
	a := NewKubeletAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"daemon_status/kubelet_status": []byte("Active: inactive (dead)"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].ENName != "KUBELET_SERVICE_STOPPED" {
		t.Errorf("ENName = %v, want %v", issues[0].ENName, "KUBELET_SERVICE_STOPPED")
	}
}

func TestNewContainerRuntimeAnalyzer(t *testing.T) {
	a := NewContainerRuntimeAnalyzer()
	if a == nil {
		t.Fatal("NewContainerRuntimeAnalyzer() returned nil")
	}
	if a.Name() != "process.container_runtime" {
		t.Errorf("Name() = %v, want %v", a.Name(), "process.container_runtime")
	}
}

func TestContainerRuntimeAnalyzerAnalyze_NoIssues(t *testing.T) {
	a := NewContainerRuntimeAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"daemon_status/docker_status":    []byte("Active: active (running)"),
		"daemon_status/containerd_status": []byte("Active: active (running)"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}

func TestParseServiceStatus(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"Active: active (running)", "running"},
		{"Active: failed (Result: exit-code)", "failed"},
		{"Active: inactive (dead)", "stopped"},
		{"Active: inactive", "stopped"},
		{"No status info", "unknown"},
	}

	for _, tt := range tests {
		result := parseServiceStatus(tt.content)
		if result != tt.expected {
			t.Errorf("parseServiceStatus(%q) = %v, want %v", tt.content, result, tt.expected)
		}
	}
}

func TestNewFirewalldAnalyzer(t *testing.T) {
	a := NewFirewalldAnalyzer()
	if a == nil {
		t.Fatal("NewFirewalldAnalyzer() returned nil")
	}
	if a.Name() != "process.firewalld" {
		t.Errorf("Name() = %v, want %v", a.Name(), "process.firewalld")
	}
}

func TestNewPIDLeakAnalyzer(t *testing.T) {
	a := NewPIDLeakAnalyzer()
	if a == nil {
		t.Fatal("NewPIDLeakAnalyzer() returned nil")
	}
	if a.Name() != "process.pid_leak" {
		t.Errorf("Name() = %v, want %v", a.Name(), "process.pid_leak")
	}
}
