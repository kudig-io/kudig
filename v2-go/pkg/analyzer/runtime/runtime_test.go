package runtime

import (
	"context"
	"testing"

	"github.com/kudig/kudig/pkg/types"
)

func TestNewDockerAnalyzer(t *testing.T) {
	a := NewDockerAnalyzer()
	if a == nil {
		t.Fatal("NewDockerAnalyzer() returned nil")
	}
	if a.Name() != "runtime.docker" {
		t.Errorf("Name() = %v, want %v", a.Name(), "runtime.docker")
	}
}

func TestDockerAnalyzerAnalyze_NoFile(t *testing.T) {
	a := NewDockerAnalyzer()
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

func TestDockerAnalyzerAnalyze_NoIssues(t *testing.T) {
	a := NewDockerAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"logs/docker.log": []byte("Docker daemon started successfully"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}

func TestDockerAnalyzerAnalyze_StartFailed(t *testing.T) {
	a := NewDockerAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"logs/docker.log": []byte("Failed to start Docker daemon"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].ENName != "DOCKER_START_FAILED" {
		t.Errorf("ENName = %v, want %v", issues[0].ENName, "DOCKER_START_FAILED")
	}
	if issues[0].Severity != types.SeverityCritical {
		t.Errorf("Severity = %v, want %v", issues[0].Severity, types.SeverityCritical)
	}
}

func TestDockerAnalyzerAnalyze_StorageDriverError(t *testing.T) {
	a := NewDockerAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"logs/docker.log": []byte("Storage driver error occurred"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].ENName != "DOCKER_STORAGE_DRIVER_ERROR" {
		t.Errorf("ENName = %v, want %v", issues[0].ENName, "DOCKER_STORAGE_DRIVER_ERROR")
	}
}

func TestNewContainerdAnalyzer(t *testing.T) {
	a := NewContainerdAnalyzer()
	if a == nil {
		t.Fatal("NewContainerdAnalyzer() returned nil")
	}
	if a.Name() != "runtime.containerd" {
		t.Errorf("Name() = %v, want %v", a.Name(), "runtime.containerd")
	}
}

func TestContainerdAnalyzerAnalyze_NoFile(t *testing.T) {
	a := NewContainerdAnalyzer()
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

func TestContainerdAnalyzerAnalyze_CreateFailed(t *testing.T) {
	a := NewContainerdAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	// Need more than 10 occurrences to trigger
	content := ""
	for i := 0; i < 12; i++ {
		content += "failed to create container "
	}
	data.RawFiles = map[string][]byte{
		"logs/containerd.log": []byte(content),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].ENName != "CONTAINER_CREATE_FAILED" {
		t.Errorf("ENName = %v, want %v", issues[0].ENName, "CONTAINER_CREATE_FAILED")
	}
}

func TestNewTimeSyncAnalyzer(t *testing.T) {
	a := NewTimeSyncAnalyzer()
	if a == nil {
		t.Fatal("NewTimeSyncAnalyzer() returned nil")
	}
	if a.Name() != "runtime.time_sync" {
		t.Errorf("Name() = %v, want %v", a.Name(), "runtime.time_sync")
	}
}

func TestNewConfigAnalyzer(t *testing.T) {
	a := NewConfigAnalyzer()
	if a == nil {
		t.Fatal("NewConfigAnalyzer() returned nil")
	}
	if a.Name() != "runtime.config" {
		t.Errorf("Name() = %v, want %v", a.Name(), "runtime.config")
	}
}
