// Package main is the entry point for kudig CLI
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/kudig/kudig/pkg/analyzer"

	// Import analyzers to register them.
	_ "github.com/kudig/kudig/pkg/analyzer/kernel"
	_ "github.com/kudig/kudig/pkg/analyzer/kubernetes"
	_ "github.com/kudig/kudig/pkg/analyzer/network"
	_ "github.com/kudig/kudig/pkg/analyzer/process"
	_ "github.com/kudig/kudig/pkg/analyzer/runtime"
	_ "github.com/kudig/kudig/pkg/analyzer/security"
	_ "github.com/kudig/kudig/pkg/analyzer/servicemesh"
	_ "github.com/kudig/kudig/pkg/analyzer/system"
	_ "github.com/kudig/kudig/pkg/ebpf/analyzer"
	"github.com/kudig/kudig/pkg/collector"
	_ "github.com/kudig/kudig/pkg/collector/offline"
	"github.com/kudig/kudig/pkg/collector/online"
	_ "github.com/kudig/kudig/pkg/collector/online"
	"github.com/kudig/kudig/pkg/history"
	"github.com/kudig/kudig/pkg/legacy"
	"github.com/kudig/kudig/pkg/metrics"
	"github.com/kudig/kudig/pkg/notifier"
	"github.com/kudig/kudig/pkg/autofix"
	"github.com/kudig/kudig/pkg/cost"
	"github.com/kudig/kudig/pkg/ai"
	"github.com/kudig/kudig/pkg/rca"
	"github.com/kudig/kudig/pkg/reporter"
	"github.com/kudig/kudig/pkg/scanner"
	"github.com/kudig/kudig/pkg/rules"
	"github.com/kudig/kudig/pkg/tui"
	"github.com/kudig/kudig/pkg/types"
)

var (
	version = "2.0.0"
	// Global flags
	verbose    bool
	outputFile string
	format     string

	// Online mode flags
	kubeconfig    string
	kubeCtx       string
	nodeName      string
	namespace     string
	allNodes      bool
	serveMetrics  bool
	metricsPort   int

	// Rules mode flags
	rulesFile string
	rulesDir  string
	listRules bool

	// RCA mode flags
	enableRCA bool

	// Pprof flags
	pprofPort int

	// Trace flags
	jaegerEndpoint string

	// Multicluster flags
	allContexts bool
	contexts    []string

	// AI flags
	aiOnline bool
)

// exitError 用于传递退出码而不直接调用 os.Exit。
type exitError struct {
	code int
}

func (e *exitError) Error() string {
	return fmt.Sprintf("exit code %d", e.code)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// 检查是否为 exitError，如果是则使用指定的退出码
		if exitErr, ok := err.(*exitError); ok {
			os.Exit(exitErr.code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "kudig",
	Short: "Kubernetes Diagnostic Toolkit",
	Long: `kudig (Kubernetes Diagnostic Toolkit) is a comprehensive diagnostic tool
for analyzing Kubernetes node issues.

It supports:
- Offline analysis of diagnostic data collected by diagnose_k8s.sh
- Online diagnosis via K8s API (real-time cluster analysis)
- Legacy mode using the original kudig.sh script`,
	Version: version,
}

var offlineCmd = &cobra.Command{
	Use:   "offline <diagnose_path>",
	Short: "Analyze diagnostic data from a directory",
	Long: `Analyze diagnostic data collected by diagnose_k8s.sh.

Example:
  kudig offline /tmp/diagnose_1702468800
  kudig offline --format json /tmp/diagnose_1702468800`,
	Args: cobra.ExactArgs(1),
	RunE: runOffline,
}

var legacyCmd = &cobra.Command{
	Use:   "legacy <diagnose_path>",
	Short: "Run analysis using legacy kudig.sh script",
	Long: `Run the original kudig.sh script for backward compatibility.

Example:
  kudig legacy /tmp/diagnose_1702468800`,
	Args: cobra.ExactArgs(1),
	RunE: runLegacy,
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze <diagnose_path>",
	Short: "Analyze diagnostic data (alias for offline)",
	Long:  `Alias for 'kudig offline'. Analyze diagnostic data from a directory.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runOffline,
}

var listAnalyzersCmd = &cobra.Command{
	Use:   "list-analyzers",
	Short: "List all available analyzers",
	RunE:  runListAnalyzers,
}

var onlineCmd = &cobra.Command{
	Use:   "online",
	Short: "Diagnose a live Kubernetes cluster",
	Long: `Perform real-time diagnosis of a Kubernetes cluster via K8s API.

Examples:
  # Use default kubeconfig
  kudig online

  # Specify kubeconfig and node
  kudig online --kubeconfig ~/.kube/config --node worker-1

  # Check all nodes
  kudig online --all-nodes

  # Focus on specific namespace
  kudig online --namespace my-app`,
	RunE: runOnline,
}

var rulesCmd = &cobra.Command{
	Use:   "rules <diagnose_path>",
	Short: "Run custom YAML rules against diagnostic data",
	Long: `Run custom diagnostic rules defined in YAML files.

Examples:
  # Run with built-in rules
  kudig rules /tmp/diagnose_1702468800

  # Use custom rules file
  kudig rules --file rules/custom.yaml /tmp/diagnose_1702468800

  # Use rules from directory
  kudig rules --dir rules/ /tmp/diagnose_1702468800

  # List available rules
  kudig rules --list`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRules,
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Manage diagnostic history",
	Long: `View and compare diagnostic history entries.

History is stored in ~/.kudig/history/ and includes:
- Timestamp of diagnosis
- Hostname
- Issues found
- Summary statistics`,
}

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List diagnostic history entries",
	Long:  `List all stored diagnostic history entries, sorted by timestamp (newest first).`,
	RunE:  runHistoryList,
}

var historyDiffCmd = &cobra.Command{
	Use:   "diff <id1> <id2>",
	Short: "Compare two diagnostic history entries",
	Long: `Compare two diagnostic history entries and show the differences.

Arguments:
  id1 - ID of the first history entry
  id2 - ID of the second history entry

Example:
  kudig history diff abc123 def456`,
	Args: cobra.ExactArgs(2),
	RunE: runHistoryDiff,
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for kudig.

To load completions:

Bash:
  $ source <(kudig completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ kudig completion bash > /etc/bash_completion.d/kudig
  # macOS:
  $ kudig completion bash > $(brew --prefix)/etc/bash_completion.d/kudig

Zsh:
  $ source <(kudig completion zsh)
  # To load completions for each session, execute once:
  $ kudig completion zsh > "${fpath[1]}/_kudig"

Fish:
  $ kudig completion fish | source
  # To load completions for each session, execute once:
  $ kudig completion fish > ~/.config/fish/completions/kudig.fish

PowerShell:
  PS> kudig completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> kudig completion powershell > kudig.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactArgs(1),
	RunE:                  runCompletion,
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start interactive TUI mode",
	Long: `Start kudig in interactive Terminal User Interface mode.

The TUI provides an intuitive menu-driven interface for:
- Running online diagnostics
- Analyzing offline data
- Viewing history
- Configuring options

Example:
  kudig tui`,
	RunE: runTUI,
}

var rcaCmd = &cobra.Command{
	Use:   "rca",
	Short: "Perform root cause analysis on diagnostic results",
	Long: `Analyze diagnostic issues and identify root causes.

The RCA engine correlates multiple symptoms to identify underlying
root causes and suggests remediation actions.

Examples:
  # Run RCA on current online diagnosis
  kudig online --rca

  # Run RCA on offline data
  kudig rca /tmp/diagnose_1702468800`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRCA,
}

var grafanaCmd = &cobra.Command{
	Use:   "grafana",
	Short: "Export Grafana dashboard JSON",
	Long: `Export a Grafana dashboard JSON for visualizing kudig metrics.

This command generates a Grafana dashboard JSON file that can be
imported into Grafana to visualize kudig diagnostic metrics.

Examples:
  # Export dashboard to file
  kudig grafana > kudig-dashboard.json

  # Export with custom output
  kudig grafana --output dashboard.json`,
	RunE: runGrafana,
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Auto-fix detected issues",
	Long: `Automatically fix detected issues where safe to do so.

This command attempts to automatically fix issues that have
safe, low-risk remediation actions available.

Examples:
  # Fix issues on current node (dry-run)
  kudig fix --dry-run

  # Actually fix issues
  kudig fix --confirm

  # Fix specific issue type
  kudig fix --type IMAGE_PULL`,
	RunE: runFix,
}

var costCmd = &cobra.Command{
	Use:   "cost",
	Short: "Analyze Kubernetes resource costs",
	Long: `Analyze and estimate Kubernetes resource costs.

This command calculates the estimated cost of running your
Kubernetes resources based on configured pricing.

Examples:
  # Analyze costs for current cluster
  kudig cost

  # Analyze with custom pricing
  kudig cost --cpu-price 0.03 --memory-price 0.005`,
	RunE: runCost,
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan container images for vulnerabilities",
	Long: `Scan container images for security vulnerabilities.

This command uses Trivy or other scanners to check images
for known CVEs and security issues.

Examples:
  # Scan an image
  kudig scan nginx:latest

  # Scan all images in cluster
  kudig scan --all-images`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScan,
}

var pprofCmd = &cobra.Command{
	Use:   "pprof",
	Short: "Run performance profiling on kudig",
	Long: `Start pprof profiling server for performance analysis.

This command starts a pprof server for CPU, memory, and goroutine
profiling of the kudig tool itself.

Examples:
  # Start pprof server on default port
  kudig pprof

  # Start on custom port
  kudig pprof --port 6060`,
	RunE: runPprof,
}

var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Start distributed tracing server",
	Long: `Start an OpenTelemetry tracing server for diagnostic tracing.

This command starts a Jaeger-compatible tracing server to
trace diagnostic operations across components.

Examples:
  # Start tracing server
  kudig trace

  # Export traces to Jaeger
  kudig trace --jaeger http://localhost:14268`,
	RunE: runTrace,
}

var multiclusterCmd = &cobra.Command{
	Use:   "multicluster",
	Short: "Diagnose multiple Kubernetes clusters",
	Aliases: []string{"mc"},
	Long: `Diagnose multiple Kubernetes clusters simultaneously.

This command runs diagnostics across multiple clusters defined
in kubeconfig or specified via flags.

Examples:
  # Diagnose all contexts in kubeconfig
  kudig multicluster --all-contexts

  # Diagnose specific contexts
  kudig multicluster --contexts prod-cluster,dr-cluster`,
	RunE: runMulticluster,
}

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI-assisted diagnostic analysis",
	Long: `Use AI/LLM to analyze diagnostic results and provide insights.

Requires one of: OpenAI API key, Qwen API key, or Ollama running locally.

Examples:
  # Analyze offline data with AI
  kudig ai /tmp/diagnose_1702468800

  # Analyze online cluster with AI
  kudig ai --online`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runAI,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Write output to file")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")

	// Online mode flags
	onlineCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	onlineCmd.Flags().StringVar(&kubeCtx, "context", "", "Kubernetes context to use")
	onlineCmd.Flags().StringVarP(&nodeName, "node", "n", "", "Node name to diagnose")
	onlineCmd.Flags().StringVar(&namespace, "namespace", "", "Namespace to focus on")
	onlineCmd.Flags().BoolVar(&allNodes, "all-nodes", false, "Diagnose all nodes")
	onlineCmd.Flags().BoolVar(&serveMetrics, "serve", false, "Start metrics server after diagnosis")
	onlineCmd.Flags().IntVar(&metricsPort, "metrics-port", 9090, "Port for metrics server")
	onlineCmd.Flags().BoolVar(&enableRCA, "rca", false, "Enable root cause analysis")

	// Rules mode flags
	rulesCmd.Flags().StringVar(&rulesFile, "file", "", "Path to rules YAML file")
	rulesCmd.Flags().StringVar(&rulesDir, "dir", "", "Path to rules directory")
	rulesCmd.Flags().BoolVar(&listRules, "list", false, "List available rules")

	// Add commands
	rootCmd.AddCommand(offlineCmd)
	rootCmd.AddCommand(onlineCmd)
	rootCmd.AddCommand(legacyCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(listAnalyzersCmd)
	rootCmd.AddCommand(rulesCmd)
	rootCmd.AddCommand(historyCmd)

	// Add history subcommands
	historyCmd.AddCommand(historyListCmd)
	historyCmd.AddCommand(historyDiffCmd)

	// Add completion command
	rootCmd.AddCommand(completionCmd)

	// Add TUI command
	rootCmd.AddCommand(tuiCmd)

	// Add RCA command
	rootCmd.AddCommand(rcaCmd)

	// Add Grafana command
	rootCmd.AddCommand(grafanaCmd)

	// Add Fix command
	rootCmd.AddCommand(fixCmd)

	// Add Cost command
	rootCmd.AddCommand(costCmd)

	// Add Scan command
	rootCmd.AddCommand(scanCmd)

	// Add Pprof command
	rootCmd.AddCommand(pprofCmd)
	pprofCmd.Flags().IntVar(&pprofPort, "port", 6060, "Port for pprof server")

	// Add Trace command
	rootCmd.AddCommand(traceCmd)
	traceCmd.Flags().StringVar(&jaegerEndpoint, "jaeger", "", "Jaeger collector endpoint")

	// Add Multicluster command
	rootCmd.AddCommand(multiclusterCmd)
	multiclusterCmd.Flags().BoolVar(&allContexts, "all-contexts", false, "Diagnose all kubeconfig contexts")
	multiclusterCmd.Flags().StringSliceVar(&contexts, "contexts", []string{}, "Comma-separated list of contexts to diagnose")

	// Add AI command
	rootCmd.AddCommand(aiCmd)
	aiCmd.Flags().BoolVar(&aiOnline, "online", false, "Use online mode instead of offline path")

	// Add deprecated flags for backward compatibility
	rootCmd.Flags().Bool("json", false, "Output JSON format (deprecated, use --format json)")
}

func runOffline(cmd *cobra.Command, args []string) error {
	diagnosePath := args[0]

	// Record start time for metrics
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		metrics.RecordDiagnosis(types.ModeOffline, status, duration)
	}()

	// Create context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Print header
	if verbose {
		fmt.Fprintf(os.Stderr, "================================================================\n")
		fmt.Fprintf(os.Stderr, "  kudig v%s - Kubernetes Diagnostic Toolkit\n", version)
		fmt.Fprintf(os.Stderr, "================================================================\n\n")
		fmt.Fprintf(os.Stderr, "诊断目录: %s\n", diagnosePath)
		fmt.Fprintf(os.Stderr, "分析时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	// Get collector
	col, ok := collector.GetCollector(types.ModeOffline)
	if !ok {
		status = "collector_error"
		return fmt.Errorf("offline collector not available")
	}

	// Collect data
	config := collector.NewOfflineConfig(diagnosePath)
	data, err := col.Collect(ctx, config)
	if err != nil {
		status = "collect_error"
		return fmt.Errorf("failed to collect data: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "节点信息: %s\n", data.NodeInfo.Hostname)
		fmt.Fprintf(os.Stderr, "\n开始诊断检查...\n\n")
	}

	// Run analyzers
	results, err := analyzer.DefaultRegistry.ExecuteAll(ctx, data)
	if err != nil {
		status = "analyze_error"
		return fmt.Errorf("failed to run analyzers: %w", err)
	}

	// Record analyzer count
	metrics.RecordAnalyzers(types.ModeOffline, len(results))

	// Collect all issues
	issues := analyzer.CollectIssues(results)

	// Record issues metrics
	metrics.RecordIssues(issues)

	// Deduplicate and sort
	issues = reporter.DeduplicateIssues(issues)
	issues = reporter.SortIssuesBySeverity(issues)

	// Generate report
	metadata := reporter.NewReportMetadata()
	metadata.Hostname = data.NodeInfo.Hostname
	metadata.DiagnosePath = diagnosePath
	metadata.Mode = "offline"

	// Determine format
	outputFormat := format
	if jsonFlag, err := cmd.Flags().GetBool("json"); err == nil && jsonFlag {
		outputFormat = "json"
	}

	rep, ok := reporter.GetReporter(outputFormat)
	if !ok {
		return fmt.Errorf("unknown format: %s", outputFormat)
	}

	output, err := rep.Generate(issues, metadata)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, output, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "报告已保存到: %s\n", outputFile)
		}
	} else {
		fmt.Println(string(output))
	}

	// Save to history
	if histMgr, err := history.NewManager(); err == nil {
		if _, err := histMgr.Save(data.NodeInfo.Hostname, "offline", issues); err == nil && verbose {
			fmt.Fprintf(os.Stderr, "已保存到历史记录\n")
		}
	}

	// Send notifications if configured
	sendNotification(data.NodeInfo.Hostname, "offline", issues)

	// Return exit code based on severity
	maxSev := types.MaxSeverity(issues)
	if maxSev == types.SeverityCritical {
		status = "critical_issues"
		return &exitError{code: 2}
	} else if len(issues) > 0 {
		status = "issues_found"
		return &exitError{code: 1}
	}

	return nil
}

func runLegacy(_ *cobra.Command, args []string) error {
	diagnosePath := args[0]

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create legacy collector
	legacyCol, err := legacy.NewLegacyCollector("")
	if err != nil {
		return fmt.Errorf("failed to initialize legacy mode: %w", err)
	}

	// Run legacy analysis
	report, err := legacyCol.GetReport(ctx, diagnosePath, verbose)
	if err != nil {
		return fmt.Errorf("legacy analysis failed: %w", err)
	}

	// Output the report
	if format == "json" {
		// Re-encode the report as JSON
		issues := legacy.ConvertBashReportToIssues(report)
		metadata := reporter.NewReportMetadata()
		metadata.Hostname = report.Hostname
		metadata.DiagnosePath = report.DiagnoseDir
		metadata.Engine = "bash"
		metadata.Mode = "legacy"

		rep, _ := reporter.GetReporter("json")
		output, err := rep.Generate(issues, metadata)
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	} else {
		// For text format, we let bash output directly
		issues := legacy.ConvertBashReportToIssues(report)
		metadata := reporter.NewReportMetadata()
		metadata.Hostname = report.Hostname
		metadata.DiagnosePath = report.DiagnoseDir
		metadata.Engine = "bash"
		metadata.Mode = "legacy"

		rep, _ := reporter.GetReporter("text")
		output, err := rep.Generate(issues, metadata)
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	}

	// Return appropriate exit code
	if report.Summary.Critical > 0 {
		return &exitError{code: 2}
	} else if report.Summary.Total > 0 {
		return &exitError{code: 1}
	}

	return nil
}

func runListAnalyzers(_ *cobra.Command, _ []string) error {
	analyzers := analyzer.DefaultRegistry.List()

	if len(analyzers) == 0 {
		fmt.Println("No analyzers registered.")
		return nil
	}

	fmt.Println("Available Analyzers:")
	fmt.Println("--------------------")

	for _, a := range analyzers {
		modes := make([]string, len(a.SupportedModes()))
		for i, m := range a.SupportedModes() {
			modes[i] = m.String()
		}

		fmt.Printf("  %s\n", a.Name())
		fmt.Printf("    Category: %s\n", a.Category())
		fmt.Printf("    Description: %s\n", a.Description())
		fmt.Printf("    Modes: %v\n", modes)
		fmt.Println()
	}

	return nil
}

func runOnline(_ *cobra.Command, _ []string) error {
	// Start metrics server if requested
	if serveMetrics {
		addr := fmt.Sprintf(":%d", metricsPort)
		metricsServer := metrics.NewServer(addr)
		go func() {
			if verbose {
				fmt.Fprintf(os.Stderr, "Starting metrics server on %s/metrics\n", addr)
			}
			if err := metricsServer.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Metrics server error: %v\n", err)
			}
		}()
	}

	// Record start time for metrics
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		metrics.RecordDiagnosis(types.ModeOnline, status, duration)
	}()

	// Create context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Print header
	if verbose {
		fmt.Fprintf(os.Stderr, "================================================================\n")
		fmt.Fprintf(os.Stderr, "  kudig v%s - Kubernetes Diagnostic Toolkit (Online Mode)\n", version)
		fmt.Fprintf(os.Stderr, "================================================================\n\n")
		fmt.Fprintf(os.Stderr, "分析时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		if nodeName != "" {
			fmt.Fprintf(os.Stderr, "目标节点: %s\n", nodeName)
		}
		if namespace != "" {
			fmt.Fprintf(os.Stderr, "目标命名空间: %s\n", namespace)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Get collector
	col, ok := collector.GetCollector(types.ModeOnline)
	if !ok {
		status = "collector_error"
		return fmt.Errorf("online collector not available")
	}

	// Handle all-nodes mode with concurrent collection
	if allNodes {
		return runOnlineAllNodes(ctx, col, status)
	}

	// Single node mode
	return runOnlineSingleNode(ctx, col, status)
}

func runOnlineSingleNode(ctx context.Context, col collector.Collector, status string) error {
	// Build config
	config := &collector.Config{
		Kubeconfig:     kubeconfig,
		Context:        kubeCtx,
		NodeName:       nodeName,
		Namespace:      namespace,
		AllNodes:       false,
		TimeoutSeconds: 60,
	}

	// Collect data
	if verbose {
		fmt.Fprintf(os.Stderr, "正在连接 Kubernetes 集群...\n")
	}

	data, err := col.Collect(ctx, config)
	if err != nil {
		status = "collect_error"
		return fmt.Errorf("failed to collect data: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "已连接到集群\n")
		if data.NodeInfo.Hostname != "" {
			fmt.Fprintf(os.Stderr, "节点: %s\n", data.NodeInfo.Hostname)
			fmt.Fprintf(os.Stderr, "Kubelet版本: %s\n", data.NodeInfo.KubeletVersion)
			fmt.Fprintf(os.Stderr, "容器运行时: %s\n", data.NodeInfo.ContainerRuntime)
		}
		fmt.Fprintf(os.Stderr, "\n开始诊断检查...\n\n")
	}

	// Run analyzers that support online mode
	results, err := analyzer.DefaultRegistry.ExecuteByMode(ctx, data, types.ModeOnline)
	if err != nil {
		status = "analyze_error"
		return fmt.Errorf("failed to run analyzers: %w", err)
	}

	// Record analyzer count
	metrics.RecordAnalyzers(types.ModeOnline, len(results))

	// Collect all issues
	issues := analyzer.CollectIssues(results)

	// Record issues metrics
	metrics.RecordIssues(issues)

	// Deduplicate and sort
	issues = reporter.DeduplicateIssues(issues)
	issues = reporter.SortIssuesBySeverity(issues)

	// Generate report
	metadata := reporter.NewReportMetadata()
	metadata.Hostname = data.NodeInfo.Hostname
	metadata.Mode = "online"
	if kubeconfig != "" {
		metadata.DiagnosePath = kubeconfig
	}

	rep, ok := reporter.GetReporter(format)
	if !ok {
		return fmt.Errorf("unknown format: %s", format)
	}

	output, err := rep.Generate(issues, metadata)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, output, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "报告已保存到: %s\n", outputFile)
		}
	} else {
		fmt.Println(string(output))
	}

	// Save to history
	if histMgr, err := history.NewManager(); err == nil {
		if _, err := histMgr.Save(data.NodeInfo.Hostname, "online", issues); err == nil && verbose {
			fmt.Fprintf(os.Stderr, "已保存到历史记录\n")
		}
	}

	// Send notifications if configured
	sendNotification(data.NodeInfo.Hostname, "online", issues)

	// Return exit code based on severity
	maxSev := types.MaxSeverity(issues)
	if maxSev == types.SeverityCritical {
		status = "critical_issues"
		return &exitError{code: 2}
	} else if len(issues) > 0 {
		status = "issues_found"
		return &exitError{code: 1}
	}

	return nil
}

func runOnlineAllNodes(ctx context.Context, col collector.Collector, status string) error {
	// Build config for collector
	config := &collector.Config{
		Kubeconfig:     kubeconfig,
		Context:        kubeCtx,
		NodeName:       "",
		Namespace:      namespace,
		AllNodes:       true,
		TimeoutSeconds: 60,
	}

	// Get online collector for concurrent collection
	onlineCollector, ok := col.(*online.Collector)
	if !ok {
		status = "collector_type_error"
		return fmt.Errorf("collector is not an online collector")
	}

	// Create progress bar
	var bar *progressbar.ProgressBar
	if verbose {
		fmt.Fprintf(os.Stderr, "正在收集所有节点数据...\n")
	} else {
		bar = progressbar.NewOptions(-1,
			progressbar.OptionSetDescription("Diagnosing nodes"),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
			progressbar.OptionSetItsString("nodes"),
			progressbar.OptionThrottle(100*time.Millisecond),
			progressbar.OptionShowElapsedTimeOnFinish(),
		)
	}

	// Collect data from all nodes concurrently
	progressFn := func(current, total int, nodeName string) {
		if bar != nil {
			bar.ChangeMax(total)
			bar.Set(current)
		} else if verbose {
			fmt.Fprintf(os.Stderr, "  [%d/%d] Diagnosing node: %s\n", current, total, nodeName)
		}
	}

	nodeResults, err := onlineCollector.CollectAllNodesConcurrent(ctx, config, progressFn)
	if err != nil {
		status = "collect_error"
		return fmt.Errorf("failed to collect nodes data: %w", err)
	}

	if bar != nil {
		bar.Finish()
		fmt.Fprintln(os.Stderr)
	}

	// Aggregate all issues from all nodes
	allIssues := make([]types.Issue, 0)
	successfulNodes := 0
	failedNodes := 0

	for _, result := range nodeResults {
		if result.Error != nil {
			failedNodes++
			fmt.Fprintf(os.Stderr, "Warning: failed to diagnose node %s: %v\n", result.NodeName, result.Error)
			continue
		}

		successfulNodes++

		// Run analyzers for this node
		results, err := analyzer.DefaultRegistry.ExecuteByMode(ctx, result.Data, types.ModeOnline)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to run analyzers on node %s: %v\n", result.NodeName, err)
			continue
		}

		// Add node name to issues for identification
		issues := analyzer.CollectIssues(results)
		for i := range issues {
			if issues[i].Metadata == nil {
				issues[i].Metadata = make(map[string]string)
			}
			issues[i].Metadata["node"] = result.NodeName
			// Prefix the issue details with node name
			issues[i].Details = fmt.Sprintf("[%s] %s", result.NodeName, issues[i].Details)
		}

		allIssues = append(allIssues, issues...)
		metrics.RecordAnalyzers(types.ModeOnline, len(results))
	}

	if successfulNodes == 0 {
		status = "all_nodes_failed"
		return fmt.Errorf("failed to diagnose any node")
	}

	// Record issues metrics
	metrics.RecordIssues(allIssues)

	// Deduplicate and sort
	allIssues = reporter.DeduplicateIssues(allIssues)
	allIssues = reporter.SortIssuesBySeverity(allIssues)

	// Generate aggregated report
	metadata := reporter.NewReportMetadata()
	metadata.Hostname = fmt.Sprintf("%d nodes", successfulNodes)
	metadata.Mode = "online-multi-node"
	if kubeconfig != "" {
		metadata.DiagnosePath = kubeconfig
	}

	rep, ok := reporter.GetReporter(format)
	if !ok {
		return fmt.Errorf("unknown format: %s", format)
	}

	output, err := rep.Generate(allIssues, metadata)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, output, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "\n报告已保存到: %s\n", outputFile)
		}
	} else {
		fmt.Println(string(output))
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "\n诊断完成: %d 个节点成功, %d 个节点失败\n", successfulNodes, failedNodes)
	}

	// Save to history
	if histMgr, err := history.NewManager(); err == nil {
		if _, err := histMgr.Save(fmt.Sprintf("%d nodes", successfulNodes), "online-multi-node", allIssues); err == nil && verbose {
			fmt.Fprintf(os.Stderr, "已保存到历史记录\n")
		}
	}

	// Send notifications if configured
	sendNotification(fmt.Sprintf("%d nodes", successfulNodes), "online-multi-node", allIssues)

	// Return exit code based on severity
	maxSev := types.MaxSeverity(allIssues)
	if maxSev == types.SeverityCritical {
		status = "critical_issues"
		return &exitError{code: 2}
	} else if len(allIssues) > 0 {
		status = "issues_found"
		return &exitError{code: 1}
	}

	return nil
}

func runRules(_ *cobra.Command, args []string) error {
	// Load rules
	loader := rules.NewLoader()

	// Load built-in rules first
	if err := loader.LoadBuiltin(); err != nil {
		return fmt.Errorf("failed to load built-in rules: %w", err)
	}

	// Load custom rules from file or directory
	if rulesFile != "" {
		if err := loader.LoadFile(rulesFile); err != nil {
			return fmt.Errorf("failed to load rules file: %w", err)
		}
	}
	if rulesDir != "" {
		if err := loader.LoadDir(rulesDir); err != nil {
			return fmt.Errorf("failed to load rules directory: %w", err)
		}
	}

	// List rules mode
	if listRules {
		allRules := loader.GetAllRules()
		fmt.Println("Available Rules:")
		fmt.Println("----------------")
		for _, r := range allRules {
			fmt.Printf("  %s\n", r.ID)
			fmt.Printf("    Name: %s\n", r.Name)
			fmt.Printf("    Category: %s\n", r.Category)
			fmt.Printf("    Severity: %s\n", r.Severity)
			fmt.Printf("    Description: %s\n", r.Description)
			fmt.Println()
		}
		return nil
	}

	// Require diagnose path for analysis
	if len(args) < 1 {
		return fmt.Errorf("diagnose_path is required for rules analysis")
	}
	diagnosePath := args[0]

	// Record start time for metrics
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		metrics.RecordDiagnosis(types.ModeOffline, status, duration)
	}()

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	if verbose {
		fmt.Fprintf(os.Stderr, "================================================================\n")
		fmt.Fprintf(os.Stderr, "  kudig v%s - Rules Engine\n", version)
		fmt.Fprintf(os.Stderr, "================================================================\n\n")
		fmt.Fprintf(os.Stderr, "诊断目录: %s\n", diagnosePath)
		fmt.Fprintf(os.Stderr, "规则数量: %d\n", len(loader.GetAllRules()))
		fmt.Fprintf(os.Stderr, "分析时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	// Collect data
	col, ok := collector.GetCollector(types.ModeOffline)
	if !ok {
		status = "collector_error"
		return fmt.Errorf("offline collector not available")
	}

	config := collector.NewOfflineConfig(diagnosePath)
	data, err := col.Collect(ctx, config)
	if err != nil {
		status = "collect_error"
		return fmt.Errorf("failed to collect data: %w", err)
	}

	// Run rule engine
	engine := rules.NewEngine(loader)
	issues, err := engine.Evaluate(ctx, data)
	if err != nil {
		status = "rules_error"
		return fmt.Errorf("failed to evaluate rules: %w", err)
	}

	// Record issues metrics
	metrics.RecordIssues(issues)

	// Deduplicate and sort
	issues = reporter.DeduplicateIssues(issues)
	issues = reporter.SortIssuesBySeverity(issues)

	// Generate report
	metadata := reporter.NewReportMetadata()
	metadata.Hostname = data.NodeInfo.Hostname
	metadata.DiagnosePath = diagnosePath
	metadata.Mode = "rules"
	metadata.Engine = "rules"

	rep, ok := reporter.GetReporter(format)
	if !ok {
		return fmt.Errorf("unknown format: %s", format)
	}

	output, err := rep.Generate(issues, metadata)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Write output
	if outputFile != "" {
		if err := os.WriteFile(outputFile, output, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "报告已保存到: %s\n", outputFile)
		}
	} else {
		fmt.Println(string(output))
	}

	// Return exit code based on severity
	maxSev := types.MaxSeverity(issues)
	if maxSev == types.SeverityCritical {
		status = "critical_issues"
		return &exitError{code: 2}
	} else if len(issues) > 0 {
		status = "issues_found"
		return &exitError{code: 1}
	}

	return nil
}

func runHistoryList(_ *cobra.Command, _ []string) error {
	mgr, err := history.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create history manager: %w", err)
	}

	entries, err := mgr.List()
	if err != nil {
		return fmt.Errorf("failed to list history: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No history entries found.")
		fmt.Println("Run a diagnosis first to create history entries.")
		return nil
	}

	fmt.Println("Diagnostic History:")
	fmt.Println("===================")
	fmt.Printf("%-18s %-20s %-15s %-10s %-10s %-10s\n", "ID", "Timestamp", "Hostname", "Critical", "Warning", "Info")
	fmt.Println(strings.Repeat("-", 90))

	for _, entry := range entries {
		shortID := entry.ID
		if len(shortID) > 16 {
			shortID = shortID[:16]
		}
		fmt.Printf("%-18s %-20s %-15s %-10d %-10d %-10d\n",
			shortID,
			entry.Timestamp.Format("2006-01-02 15:04"),
			truncate(entry.Hostname, 15),
			entry.Summary.Critical,
			entry.Summary.Warning,
			entry.Summary.Info,
		)
	}

	fmt.Printf("\nTotal entries: %d\n", len(entries))
	fmt.Println("\nUse 'kudig history diff <id1> <id2>' to compare entries.")

	return nil
}

func runHistoryDiff(_ *cobra.Command, args []string) error {
	id1, id2 := args[0], args[1]

	mgr, err := history.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create history manager: %w", err)
	}

	diff, err := mgr.Diff(id1, id2)
	if err != nil {
		return fmt.Errorf("failed to diff history entries: %w", err)
	}

	fmt.Println("History Comparison:")
	fmt.Println("===================")
	fmt.Printf("Entry 1: %s (%s on %s)\n", diff.Entry1.ID[:16], diff.Entry1.Mode, diff.Entry1.Hostname)
	fmt.Printf("         %s\n", diff.Entry1.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Entry 2: %s (%s on %s)\n", diff.Entry2.ID[:16], diff.Entry2.Mode, diff.Entry2.Hostname)
	fmt.Printf("         %s\n", diff.Entry2.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Print summary changes
	fmt.Println("Summary Changes:")
	fmt.Printf("  Critical: %d → %d (%+d)\n", diff.Entry1.Summary.Critical, diff.Entry2.Summary.Critical, diff.Entry2.Summary.Critical-diff.Entry1.Summary.Critical)
	fmt.Printf("  Warning:  %d → %d (%+d)\n", diff.Entry1.Summary.Warning, diff.Entry2.Summary.Warning, diff.Entry2.Summary.Warning-diff.Entry1.Summary.Warning)
	fmt.Printf("  Info:     %d → %d (%+d)\n", diff.Entry1.Summary.Info, diff.Entry2.Summary.Info, diff.Entry2.Summary.Info-diff.Entry1.Summary.Info)
	fmt.Println()

	// Print added issues
	if len(diff.AddedIssues) > 0 {
		fmt.Printf("🆕 Added Issues (%d):\n", len(diff.AddedIssues))
		fmt.Println(strings.Repeat("-", 40))
		for _, issue := range diff.AddedIssues {
			fmt.Printf("  [%s] %s\n", severityString(issue.Severity), issue.CNName)
			fmt.Printf("      %s\n", issue.Details)
			fmt.Println()
		}
	}

	// Print removed issues
	if len(diff.RemovedIssues) > 0 {
		fmt.Printf("✅ Resolved Issues (%d):\n", len(diff.RemovedIssues))
		fmt.Println(strings.Repeat("-", 40))
		for _, issue := range diff.RemovedIssues {
			fmt.Printf("  [%s] %s\n", severityString(issue.Severity), issue.CNName)
			fmt.Printf("      %s\n", issue.Details)
			fmt.Println()
		}
	}

	if len(diff.AddedIssues) == 0 && len(diff.RemovedIssues) == 0 {
		fmt.Println("No changes detected between the two entries.")
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func severityString(s types.Severity) string {
	switch s {
	case types.SeverityCritical:
		return "CRITICAL"
	case types.SeverityWarning:
		return "WARNING"
	case types.SeverityInfo:
		return "INFO"
	default:
		return "UNKNOWN"
	}
}

// sendNotification sends webhook notifications if configured and issues meet severity threshold
func sendNotification(hostname, mode string, issues []types.Issue) {
	notifyConfig := notifier.NewConfigFromEnv()
	if !notifyConfig.ShouldNotify(issues) {
		return
	}

	multiNotifier := notifier.NewMultiNotifier(notifyConfig)
	if multiNotifier == nil || len(multiNotifier.Notifiers) == 0 {
		return
	}

	title := fmt.Sprintf("🚨 Kudig Alert: Issues detected on %s", hostname)
	message := fmt.Sprintf("Diagnostic mode: %s\nFound %d issues (%d critical, %d warning, %d info)",
		mode,
		len(issues),
		countBySeverity(issues, types.SeverityCritical),
		countBySeverity(issues, types.SeverityWarning),
		countBySeverity(issues, types.SeverityInfo),
	)

	errors := multiNotifier.Send(title, message, issues)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "Notification error: %v\n", err)
		}
	}
}

func countBySeverity(issues []types.Issue, sev types.Severity) int {
	count := 0
	for _, issue := range issues {
		if issue.Severity == sev {
			count++
		}
	}
	return count
}

func runCompletion(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		return cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return fmt.Errorf("unsupported shell type: %s. Use bash, zsh, fish, or powershell", args[0])
	}
}

func runTUI(_ *cobra.Command, _ []string) error {
	return tui.RunTUI()
}

func runRCA(_ *cobra.Command, args []string) error {
	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Collect data
	var issues []types.Issue

	if len(args) > 0 {
		// Offline mode
		diagnosePath := args[0]
		col, ok := collector.GetCollector(types.ModeOffline)
		if !ok {
			return fmt.Errorf("offline collector not available")
		}
		config := collector.NewOfflineConfig(diagnosePath)
		data, err := col.Collect(ctx, config)
		if err != nil {
			return fmt.Errorf("failed to collect data: %w", err)
		}

		// Run analyzers
		results, err := analyzer.DefaultRegistry.ExecuteAll(ctx, data)
		if err != nil {
			return fmt.Errorf("failed to run analyzers: %w", err)
		}
		issues = analyzer.CollectIssues(results)
	} else {
		// Online mode
		col, ok := collector.GetCollector(types.ModeOnline)
		if !ok {
			return fmt.Errorf("online collector not available")
		}
		config := &collector.Config{
			Kubeconfig:     kubeconfig,
			Context:        kubeCtx,
			NodeName:       nodeName,
			Namespace:      namespace,
			TimeoutSeconds: 60,
		}
		data, err := col.Collect(ctx, config)
		if err != nil {
			return fmt.Errorf("failed to collect data: %w", err)
		}

		// Run analyzers
		results, err := analyzer.DefaultRegistry.ExecuteByMode(ctx, data, types.ModeOnline)
		if err != nil {
			return fmt.Errorf("failed to run analyzers: %w", err)
		}
		issues = analyzer.CollectIssues(results)
	}

	// Perform RCA
	engine := rca.NewEngine()
	rootCauses := engine.Analyze(ctx, issues)

	// Output results
	fmt.Printf("诊断发现 %d 个问题\n\n", len(issues))
	fmt.Println(rca.FormatRootCauses(rootCauses))

	return nil
}

func runGrafana(_ *cobra.Command, _ []string) error {
	generator := reporter.NewGrafanaDashboardGenerator()
	dashboard, err := generator.GenerateDashboard()
	if err != nil {
		return fmt.Errorf("failed to generate dashboard: %w", err)
	}
	fmt.Println(string(dashboard))
	return nil
}

func runFix(_ *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var issues []types.Issue

	if len(args) > 0 {
		diagnosePath := args[0]
		col, ok := collector.GetCollector(types.ModeOffline)
		if !ok {
			return fmt.Errorf("offline collector not available")
		}
		config := collector.NewOfflineConfig(diagnosePath)
		data, err := col.Collect(ctx, config)
		if err != nil {
			return fmt.Errorf("failed to collect data: %w", err)
		}
		results, err := analyzer.DefaultRegistry.ExecuteAll(ctx, data)
		if err != nil {
			return fmt.Errorf("failed to run analyzers: %w", err)
		}
		issues = analyzer.CollectIssues(results)
	} else {
		col, ok := collector.GetCollector(types.ModeOnline)
		if !ok {
			return fmt.Errorf("online collector not available")
		}
		config := &collector.Config{
			Kubeconfig:     kubeconfig,
			Context:        kubeCtx,
			NodeName:       nodeName,
			Namespace:      namespace,
			TimeoutSeconds: 60,
		}
		data, err := col.Collect(ctx, config)
		if err != nil {
			return fmt.Errorf("failed to collect data: %w", err)
		}
		results, err := analyzer.DefaultRegistry.ExecuteAll(ctx, data)
		if err != nil {
			return fmt.Errorf("failed to run analyzers: %w", err)
		}
		issues = analyzer.CollectIssues(results)
	}

	engine := autofix.NewEngine(true)

	fixable := engine.GetFixableIssues(issues)
	if len(fixable) == 0 {
		fmt.Println("没有可自动修复的问题")
		return nil
	}

	fmt.Printf("发现 %d 个可自动修复的问题（共 %d 个）:\n\n", len(fixable), len(issues))
	for _, issue := range fixable {
		action, _ := engine.CanFix(issue)
		fmt.Printf("  [%s] %s\n", issue.ENName, issue.CNName)
		fmt.Printf("    修复操作: %s (风险: %s)\n", action.Description, action.Risk)
		fmt.Printf("    命令: %s\n\n", action.Command)
	}

	fmt.Println("以上为 dry-run 模式预览。使用 --confirm 执行实际修复。")
	fmt.Println(autofix.FormatResults(engine.FixAll(ctx, fixable)))
	return nil
}

func runCost(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	col, ok := collector.GetCollector(types.ModeOnline)
	if !ok {
		return fmt.Errorf("online collector not available, cost analysis requires cluster access")
	}

	config := &collector.Config{
		Kubeconfig:     kubeconfig,
		Context:        kubeCtx,
		NodeName:       nodeName,
		Namespace:      namespace,
		TimeoutSeconds: 60,
	}

	data, err := col.Collect(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to collect data: %w", err)
	}

	a := cost.NewCostAnalyzer()
	result, err := a.Analyze(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to analyze costs: %w", err)
	}

	fmt.Println(cost.FormatResult(result))
	return nil
}

func runScan(_ *cobra.Command, args []string) error {
	image := "nginx:latest"
	if len(args) > 0 {
		image = args[0]
	}

	s := scanner.NewImageScanner()

	if !s.IsAvailable() {
		return fmt.Errorf("scanner %q not found in PATH; install trivy first: https://trivy.dev", s.ScannerType)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := s.ScanImage(ctx, image)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	fmt.Println(scanner.FormatResult(result))
	return nil
}

func runPprof(_ *cobra.Command, _ []string) error {
	addr := fmt.Sprintf("localhost:%d", pprofPort)
	fmt.Printf("Starting pprof server on http://%s/debug/pprof/\n", addr)
	fmt.Println("\nAvailable endpoints:")
	fmt.Printf("  http://%s/debug/pprof/           - Index\n", addr)
	fmt.Printf("  http://%s/debug/pprof/profile    - CPU Profile\n", addr)
	fmt.Printf("  http://%s/debug/pprof/heap       - Heap Profile\n", addr)
	fmt.Printf("  http://%s/debug/pprof/goroutine  - Goroutine Profile\n", addr)
	fmt.Printf("  http://%s/debug/pprof/allocs     - Allocations\n", addr)
	fmt.Println("\nPress Ctrl+C to stop")

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return http.ListenAndServe(addr, mux)
}

func runTrace(_ *cobra.Command, _ []string) error {
	return fmt.Errorf("trace 功能尚未实现 (experimental): OpenTelemetry 集成正在开发中")
}

func runMulticluster(_ *cobra.Command, _ []string) error {
	return fmt.Errorf("multicluster 功能尚未实现 (experimental): 多集群诊断正在开发中")
}

func runAI(_ *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	config := ai.LoadConfig()
	if config.APIKey == "" {
		return fmt.Errorf("AI 功能需要配置 API Key。请设置环境变量 KUDIG_AI_API_KEY，或使用 Ollama 本地模式")
	}

	assistant, err := ai.NewAssistant(config)
	if err != nil {
		return fmt.Errorf("failed to create AI assistant: %w", err)
	}

	if !assistant.IsAvailable() {
		return fmt.Errorf("AI provider not available (provider: %s)", config.Provider)
	}

	var issues []types.Issue

	if aiOnline {
		col, ok := collector.GetCollector(types.ModeOnline)
		if !ok {
			return fmt.Errorf("online collector not available")
		}
		colConfig := &collector.Config{
			Kubeconfig:     kubeconfig,
			Context:        kubeCtx,
			NodeName:       nodeName,
			Namespace:      namespace,
			TimeoutSeconds: 60,
		}
		data, err := col.Collect(ctx, colConfig)
		if err != nil {
			return fmt.Errorf("failed to collect data: %w", err)
		}
		results, err := analyzer.DefaultRegistry.ExecuteAll(ctx, data)
		if err != nil {
			return fmt.Errorf("failed to run analyzers: %w", err)
		}
		issues = analyzer.CollectIssues(results)
	} else if len(args) > 0 {
		col, ok := collector.GetCollector(types.ModeOffline)
		if !ok {
			return fmt.Errorf("offline collector not available")
		}
		colConfig := collector.NewOfflineConfig(args[0])
		data, err := col.Collect(ctx, colConfig)
		if err != nil {
			return fmt.Errorf("failed to collect data: %w", err)
		}
		results, err := analyzer.DefaultRegistry.ExecuteAll(ctx, data)
		if err != nil {
			return fmt.Errorf("failed to run analyzers: %w", err)
		}
		issues = analyzer.CollectIssues(results)
	} else {
		return fmt.Errorf("请指定离线诊断数据路径，或使用 --online 进行在线分析")
	}

	fmt.Printf("发现 %d 个问题，正在使用 AI 分析...\n\n", len(issues))

	result, err := assistant.AnalyzeWithAI(ctx, issues, "")
	if err != nil {
		return fmt.Errorf("AI analysis failed: %w", err)
	}

	fmt.Printf("=== AI 分析结果 ===\n\n")
	fmt.Printf("摘要: %s\n\n", result.Summary)
	if result.RootCause != "" {
		fmt.Printf("根因分析:\n%s\n\n", result.RootCause)
	}
	if len(result.Suggestions) > 0 {
		fmt.Println("修复建议:")
		for i, s := range result.Suggestions {
			fmt.Printf("  %d. [%s] %s\n", i+1, s.Risk, s.Description)
			if s.Command != "" {
				fmt.Printf("     命令: %s\n", s.Command)
			}
		}
	}
	fmt.Printf("\n置信度: %.0f%%\n", result.Confidence*100)

	return nil
}
