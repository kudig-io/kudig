// Package main implements the kubectl-kudig plugin
// This binary is named kubectl-kudig to be recognized as a kubectl plugin
// Usage: kubectl kudig [flags]
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/kudig/kudig/pkg/analyzer"

	// Import analyzers to register them
	_ "github.com/kudig/kudig/pkg/analyzer/kernel"
	_ "github.com/kudig/kudig/pkg/analyzer/kubernetes"
	_ "github.com/kudig/kudig/pkg/analyzer/network"
	_ "github.com/kudig/kudig/pkg/analyzer/process"
	_ "github.com/kudig/kudig/pkg/analyzer/runtime"
	_ "github.com/kudig/kudig/pkg/analyzer/system"
	"github.com/kudig/kudig/pkg/collector"
	_ "github.com/kudig/kudig/pkg/collector/online"
	"github.com/kudig/kudig/pkg/collector/online"
	"github.com/kudig/kudig/pkg/metrics"
	"github.com/kudig/kudig/pkg/reporter"
	"github.com/kudig/kudig/pkg/types"
)

var (
	version = "2.0.0"

	// Global flags
	verbose    bool
	outputFile string
	format     string

	// Mode flags
	allNodes   bool
	nodeName   string
	namespace  string

	// Metrics flags
	serveMetrics bool
	metricsPort  int
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "kubectl-kudig",
	Short: "Diagnose Kubernetes cluster issues",
	Long: `kudig (Kubernetes Diagnostic Toolkit) is a comprehensive diagnostic tool
for analyzing Kubernetes cluster and node issues.

This is a kubectl plugin. Use as: kubectl kudig [flags]

Examples:
  # Diagnose current node
  kubectl kudig

  # Diagnose all nodes
  kubectl kudig --all-nodes

  # Diagnose specific node
  kubectl kudig --node worker-1

  # Output as JSON
  kubectl kudig --format json

  # Generate HTML report
  kubectl kudig --format html -o report.html`,
	Version: version,
	RunE:    runDiagnose,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Write output to file")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Output format (text, json, html)")

	// Mode flags
	rootCmd.Flags().BoolVar(&allNodes, "all-nodes", false, "Diagnose all nodes")
	rootCmd.Flags().StringVarP(&nodeName, "node", "n", "", "Node name to diagnose")
	rootCmd.Flags().StringVar(&namespace, "namespace", "", "Namespace to focus on")

	// Metrics flags
	rootCmd.Flags().BoolVar(&serveMetrics, "serve", false, "Start metrics server")
	rootCmd.Flags().IntVar(&metricsPort, "metrics-port", 9090, "Port for metrics server")
}

func runDiagnose(_ *cobra.Command, _ []string) error {
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
		fmt.Fprintf(os.Stderr, "  kubectl kudig v%s - Kubernetes Diagnostic Toolkit\n", version)
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
		return runAllNodes(ctx, col, status)
	}

	// Single node mode
	return runSingleNode(ctx, col, status)
}

func runSingleNode(ctx context.Context, col collector.Collector, status string) error {
	// Build config
	config := &collector.Config{
		Kubeconfig:     "",
		Context:        "",
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
		os.Exit(2)
	} else if len(issues) > 0 {
		status = "issues_found"
		os.Exit(1)
	}

	return nil
}

func runAllNodes(ctx context.Context, col collector.Collector, status string) error {
	// Build config for collector
	config := &collector.Config{
		Kubeconfig:     "",
		Context:        "",
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

	// Return exit code based on severity
	maxSev := types.MaxSeverity(allIssues)
	if maxSev == types.SeverityCritical {
		status = "critical_issues"
		os.Exit(2)
	} else if len(allIssues) > 0 {
		status = "issues_found"
		os.Exit(1)
	}

	return nil
}
