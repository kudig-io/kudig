package analyzer

import (
	"context"
	"testing"

	"github.com/kudig/kudig/pkg/types"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if len(r.List()) != 0 {
		t.Error("New registry should be empty")
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()
	mock := newMockAnalyzer()

	err := r.Register(mock)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if len(r.List()) != 1 {
		t.Errorf("List() length = %v, want %v", len(r.List()), 1)
	}
}

func TestRegistryRegisterDuplicate(t *testing.T) {
	r := NewRegistry()
	mock := newMockAnalyzer()

	if err := r.Register(mock); err != nil {
		t.Fatalf("First Register() error = %v", err)
	}

	// Try to register again
	err := r.Register(mock)
	if err == nil {
		t.Error("Register() should return error for duplicate")
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry()
	mock := newMockAnalyzer()
	r.Register(mock)

	got, ok := r.Get("mock.analyzer")
	if !ok {
		t.Error("Get() should find registered analyzer")
	}
	if got.Name() != "mock.analyzer" {
		t.Errorf("Name = %v, want %v", got.Name(), "mock.analyzer")
	}

	_, ok = r.Get("nonexistent")
	if ok {
		t.Error("Get() should not find nonexistent analyzer")
	}
}

func TestRegistryList(t *testing.T) {
	r := NewRegistry()

	// Register multiple analyzers
	r.Register(newMockAnalyzer())
	r.Register(&mockAnalyzer2{})

	list := r.List()
	if len(list) != 2 {
		t.Errorf("List() length = %v, want %v", len(list), 2)
	}
}

func TestRegistryListByCategory(t *testing.T) {
	r := NewRegistry()

	// Register analyzers with different categories
	r.Register(newMockAnalyzer()) // category: "test"
	r.Register(&mockAnalyzer2{})  // category: "test2"

	testAnalyzers := r.ListByCategory("test")
	if len(testAnalyzers) != 1 {
		t.Errorf("ListByCategory('test') length = %v, want %v", len(testAnalyzers), 1)
	}

	test2Analyzers := r.ListByCategory("test2")
	if len(test2Analyzers) != 1 {
		t.Errorf("ListByCategory('test2') length = %v, want %v", len(test2Analyzers), 1)
	}

	noneAnalyzers := r.ListByCategory("nonexistent")
	if len(noneAnalyzers) != 0 {
		t.Errorf("ListByCategory('nonexistent') length = %v, want %v", len(noneAnalyzers), 0)
	}
}

func TestRegistryListByMode(t *testing.T) {
	r := NewRegistry()
	r.Register(newMockAnalyzer()) // supports both modes
	r.Register(&mockAnalyzer2{})  // supports only offline

	offlineAnalyzers := r.ListByMode(types.ModeOffline)
	if len(offlineAnalyzers) != 2 {
		t.Errorf("ListByMode(Offline) length = %v, want %v", len(offlineAnalyzers), 2)
	}

	onlineAnalyzers := r.ListByMode(types.ModeOnline)
	if len(onlineAnalyzers) != 1 {
		t.Errorf("ListByMode(Online) length = %v, want %v", len(onlineAnalyzers), 1)
	}
}

func TestCollectIssues(t *testing.T) {
	results := []Result{
		{
			AnalyzerName: "analyzer1",
			Issues: []types.Issue{
				{Severity: types.SeverityCritical, ENName: "ISSUE1"},
			},
		},
		{
			AnalyzerName: "analyzer2",
			Issues: []types.Issue{
				{Severity: types.SeverityWarning, ENName: "ISSUE2"},
				{Severity: types.SeverityInfo, ENName: "ISSUE3"},
			},
		},
	}

	issues := CollectIssues(results)
	if len(issues) != 3 {
		t.Errorf("CollectIssues length = %v, want %v", len(issues), 3)
	}
}

func TestCollectIssuesEmpty(t *testing.T) {
	results := []Result{}
	issues := CollectIssues(results)
	if len(issues) != 0 {
		t.Errorf("CollectIssues length = %v, want %v", len(issues), 0)
	}
}

// mockAnalyzer2 is another test double with different category
type mockAnalyzer2 struct {
	BaseAnalyzer
}

func (m *mockAnalyzer2) Name() string        { return "mock.analyzer2" }
func (m *mockAnalyzer2) Description() string { return "Mock analyzer 2" }
func (m *mockAnalyzer2) Category() string    { return "test2" }
func (m *mockAnalyzer2) Analyze(_ context.Context, _ *types.DiagnosticData) ([]types.Issue, error) {
	return nil, nil
}
func (m *mockAnalyzer2) SupportedModes() []types.DataMode { return []types.DataMode{types.ModeOffline} }
func (m *mockAnalyzer2) Dependencies() []string           { return nil }

func (m *mockAnalyzer2) initBase() {
	m.BaseAnalyzer = *NewBaseAnalyzer(
		"mock.analyzer2",
		"Mock analyzer 2",
		"test2",
		[]types.DataMode{types.ModeOffline},
	)
}
