package tui

import (
	"strings"
	"testing"

	"github.com/kudig/kudig/pkg/types"
)

func containsEmoji(title string, severity types.Severity) bool {
	switch severity {
	case types.SeverityCritical:
		return strings.Contains(title, "🔴")
	case types.SeverityWarning:
		return strings.Contains(title, "🟡")
	case types.SeverityInfo:
		return strings.Contains(title, "🟢")
	}
	return false
}

func TestNewModel(t *testing.T) {
	m := NewModel()
	if m.state != StateMenu {
		t.Errorf("Expected initial state StateMenu, got %d", m.state)
	}
	if m.onlineMode != true {
		t.Error("Expected onlineMode to be true by default")
	}
}

func TestModel_Init(t *testing.T) {
	m := NewModel()
	cmd := m.Init()
	if cmd != nil {
		t.Error("Expected nil Init command")
	}
}

func TestMenuItem_FilterValue(t *testing.T) {
	item := MenuItem{Title: "Test Item"}
	if item.FilterValue() != "Test Item" {
		t.Errorf("Expected FilterValue 'Test Item', got '%s'", item.FilterValue())
	}
}

func TestIssueItem_FilterValue(t *testing.T) {
	item := IssueItem{Issue: types.Issue{CNName: "CPU过高", ENName: "HIGH_CPU"}}
	expected := "CPU过高 HIGH_CPU"
	if item.FilterValue() != expected {
		t.Errorf("Expected FilterValue '%s', got '%s'", expected, item.FilterValue())
	}
}

func TestIssueItem_Title(t *testing.T) {
	tests := []struct {
		severity types.Severity
		expected string
	}{
		{types.SeverityCritical, "严重问题"},
		{types.SeverityWarning, "警告问题"},
		{types.SeverityInfo, "信息问题"},
	}
	for _, tc := range tests {
		item := IssueItem{Issue: types.Issue{Severity: tc.severity, CNName: tc.expected}}
		title := item.Title()
		if !containsEmoji(title, tc.severity) || !strings.Contains(title, tc.expected) {
			t.Errorf("Title should contain severity icon and '%s', got '%s'", tc.expected, title)
		}
	}
}

func TestIssueItem_Description(t *testing.T) {
	item := IssueItem{Issue: types.Issue{
		ENName:  "HIGH_CPU",
		Details: "CPU usage is 95%",
	}}
	desc := item.Description()
	if desc != "[HIGH_CPU] CPU usage is 95%" {
		t.Errorf("Unexpected description: %s", desc)
	}
}

func TestIssueItem_Description_Truncate(t *testing.T) {
	longDetails := ""
	for i := 0; i < 100; i++ {
		longDetails += "x"
	}
	item := IssueItem{Issue: types.Issue{
		ENName:  "TEST",
		Details: longDetails,
	}}
	desc := item.Description()
	if len(desc) > 60 {
		t.Errorf("Description should be truncated, got length %d: %s", len(desc), desc)
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"a very long string that exceeds max", 10, "a very ..."},
		{"exact", 5, "exact"},
	}
	for _, tc := range tests {
		result := truncate(tc.input, tc.maxLen)
		if result != tc.expected {
			t.Errorf("truncate(%q, %d) = %q, expected %q", tc.input, tc.maxLen, result, tc.expected)
		}
	}
}

func TestCountIssuesBySeverity(t *testing.T) {
	issues := []types.Issue{
		{Severity: types.SeverityCritical},
		{Severity: types.SeverityCritical},
		{Severity: types.SeverityWarning},
		{Severity: types.SeverityInfo},
		{Severity: types.SeverityInfo},
		{Severity: types.SeverityInfo},
	}
	critical, warning, info := countIssuesBySeverity(issues)
	if critical != 2 {
		t.Errorf("Expected 2 critical, got %d", critical)
	}
	if warning != 1 {
		t.Errorf("Expected 1 warning, got %d", warning)
	}
	if info != 3 {
		t.Errorf("Expected 3 info, got %d", info)
	}
}

func TestCountIssuesBySeverity_Empty(t *testing.T) {
	critical, warning, info := countIssuesBySeverity(nil)
	if critical != 0 || warning != 0 || info != 0 {
		t.Errorf("Expected all zeros for nil issues, got %d/%d/%d", critical, warning, info)
	}
}
