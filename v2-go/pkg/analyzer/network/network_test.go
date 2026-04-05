package network

import (
	"context"
	"testing"

	"github.com/kudig/kudig/pkg/types"
)

func TestNewInterfaceAnalyzer(t *testing.T) {
	a := NewInterfaceAnalyzer()
	if a == nil {
		t.Fatal("NewInterfaceAnalyzer() returned nil")
	}
	if a.Name() != "network.interface" {
		t.Errorf("Name() = %v, want %v", a.Name(), "network.interface")
	}
}

func TestInterfaceAnalyzerAnalyze_NoFile(t *testing.T) {
	a := NewInterfaceAnalyzer()
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

func TestInterfaceAnalyzerAnalyze_NoIssue(t *testing.T) {
	a := NewInterfaceAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"network_info": []byte("eth0: <UP,BROADCAST> mtu 1500 state UP"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}

func TestInterfaceAnalyzerAnalyze_WithDownInterface(t *testing.T) {
	a := NewInterfaceAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	data.RawFiles = map[string][]byte{
		"network_info": []byte("eth1: <BROADCAST> mtu 1500 state DOWN"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}
	if issues[0].ENName != "NETWORK_INTERFACE_DOWN" {
		t.Errorf("ENName = %v, want %v", issues[0].ENName, "NETWORK_INTERFACE_DOWN")
	}
}

func TestInterfaceAnalyzerAnalyze_ExcludeVeth(t *testing.T) {
	a := NewInterfaceAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	// veth interfaces should be excluded
	data.RawFiles = map[string][]byte{
		"network_info": []byte("veth123: <BROADCAST> mtu 1500 state DOWN"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	// veth interfaces should be excluded
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues for veth, got %d", len(issues))
	}
}

func TestInterfaceAnalyzerAnalyze_ExcludeLo(t *testing.T) {
	a := NewInterfaceAnalyzer()
	ctx := context.Background()
	data := types.NewDiagnosticData(types.ModeOffline)
	// lo interface should be excluded
	data.RawFiles = map[string][]byte{
		"network_info": []byte("lo: <BROADCAST> mtu 1500 state DOWN"),
	}

	issues, err := a.Analyze(ctx, data)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}
	// lo should be excluded
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues for lo, got %d", len(issues))
	}
}

func TestNewRouteAnalyzer(t *testing.T) {
	a := NewRouteAnalyzer()
	if a == nil {
		t.Fatal("NewRouteAnalyzer() returned nil")
	}
	if a.Name() != "network.route" {
		t.Errorf("Name() = %v, want %v", a.Name(), "network.route")
	}
}

func TestNewIptablesAnalyzer(t *testing.T) {
	a := NewIptablesAnalyzer()
	if a == nil {
		t.Fatal("NewIptablesAnalyzer() returned nil")
	}
	if a.Name() != "network.iptables" {
		t.Errorf("Name() = %v, want %v", a.Name(), "network.iptables")
	}
}

func TestNewInodeAnalyzer(t *testing.T) {
	a := NewInodeAnalyzer()
	if a == nil {
		t.Fatal("NewInodeAnalyzer() returned nil")
	}
	if a.Name() != "network.inode" {
		t.Errorf("Name() = %v, want %v", a.Name(), "network.inode")
	}
}
