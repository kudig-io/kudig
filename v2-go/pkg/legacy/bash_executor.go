// Package legacy provides backward compatibility with kudig.sh
package legacy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kudig/kudig/pkg/types"
)

// BashExecutor executes the legacy kudig.sh script
type BashExecutor struct {
	// ScriptPath is the path to kudig.sh
	ScriptPath string
}

// NewBashExecutor creates a new BashExecutor
func NewBashExecutor(scriptPath string) *BashExecutor {
	return &BashExecutor{
		ScriptPath: scriptPath,
	}
}

// FindScript attempts to locate kudig.sh in common locations
func FindScript() (string, error) {
	// Check common locations
	locations := []string{
		"./kudig.sh",
		"./scripts/kudig.sh",
		"/usr/local/bin/kudig.sh",
		"/opt/kudig/kudig.sh",
	}

	// Also check relative to executable
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		locations = append(locations,
			filepath.Join(execDir, "kudig.sh"),
			filepath.Join(execDir, "scripts", "kudig.sh"),
			filepath.Join(execDir, "..", "scripts", "kudig.sh"),
		)
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return filepath.Abs(loc)
		}
	}

	return "", fmt.Errorf("kudig.sh not found in common locations")
}

// Execute runs kudig.sh and returns the parsed issues
func (e *BashExecutor) Execute(ctx context.Context, diagnosePath string, verbose bool) ([]types.Issue, *BashReport, error) {
	// Build command arguments
	args := []string{e.ScriptPath, "--json"}
	if verbose {
		args = append(args, "--verbose")
	}
	args = append(args, diagnosePath)

	// Create command with context
	cmd := exec.CommandContext(ctx, "bash", args...)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()

	// Parse the JSON output even if there was an error (exit code 1 or 2 is expected)
	if stdout.Len() == 0 {
		if err != nil {
			return nil, nil, fmt.Errorf("bash script failed: %w, stderr: %s", err, stderr.String())
		}
		return nil, nil, fmt.Errorf("bash script produced no output")
	}

	// Parse JSON output
	var report BashReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		return nil, nil, fmt.Errorf("failed to parse bash output: %w, output: %s", err, stdout.String())
	}

	// Convert to Go Issue types
	issues := ConvertBashReportToIssues(&report)

	return issues, &report, nil
}

// BashReport matches the JSON output format of kudig.sh
type BashReport struct {
	ReportVersion string        `json:"report_version"`
	Timestamp     string        `json:"timestamp"`
	Hostname      string        `json:"hostname"`
	DiagnoseDir   string        `json:"diagnose_dir"`
	Anomalies     []BashAnomaly `json:"anomalies"`
	Summary       BashSummary   `json:"summary"`
}

// BashAnomaly represents an anomaly from kudig.sh output
type BashAnomaly struct {
	Severity string `json:"severity"`
	CNName   string `json:"cn_name"`
	ENName   string `json:"en_name"`
	Details  string `json:"details"`
	Location string `json:"location"`
}

// BashSummary represents the summary from kudig.sh output
type BashSummary struct {
	Critical int `json:"critical"`
	Warning  int `json:"warning"`
	Info     int `json:"info"`
	Total    int `json:"total"`
}

// ConvertBashReportToIssues converts BashReport to []types.Issue
func ConvertBashReportToIssues(report *BashReport) []types.Issue {
	issues := make([]types.Issue, 0, len(report.Anomalies))

	for _, anomaly := range report.Anomalies {
		issue := types.Issue{
			Severity:     types.ParseSeverity(anomaly.Severity),
			CNName:       anomaly.CNName,
			ENName:       anomaly.ENName,
			Details:      anomaly.Details,
			Location:     anomaly.Location,
			Timestamp:    time.Now(),
			AnalyzerName: "legacy.bash",
			Metadata: map[string]string{
				"source": "kudig.sh",
			},
		}
		issues = append(issues, issue)
	}

	return issues
}

// LegacyCollector implements the Collector interface using kudig.sh
type LegacyCollector struct {
	executor *BashExecutor
}

// NewLegacyCollector creates a new legacy collector
func NewLegacyCollector(scriptPath string) (*LegacyCollector, error) {
	if scriptPath == "" {
		var err error
		scriptPath, err = FindScript()
		if err != nil {
			return nil, err
		}
	}

	return &LegacyCollector{
		executor: NewBashExecutor(scriptPath),
	}, nil
}

// Execute runs the legacy analysis and returns issues
func (c *LegacyCollector) Execute(ctx context.Context, diagnosePath string, verbose bool) ([]types.Issue, error) {
	issues, _, err := c.executor.Execute(ctx, diagnosePath, verbose)
	return issues, err
}

// GetReport runs the legacy analysis and returns the full report
func (c *LegacyCollector) GetReport(ctx context.Context, diagnosePath string, verbose bool) (*BashReport, error) {
	_, report, err := c.executor.Execute(ctx, diagnosePath, verbose)
	return report, err
}
