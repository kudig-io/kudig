// Package main is the entry point for kudig CLI
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/kudig/kudig/pkg/analyzer"
	"github.com/kudig/kudig/pkg/collector"
	_ "github.com/kudig/kudig/pkg/collector/offline"
	_ "github.com/kudig/kudig/pkg/collector/online"
	"github.com/kudig/kudig/pkg/legacy"
	"github.com/kudig/kudig/pkg/reporter"
	"github.com/kudig/kudig/pkg/rules"
	"github.com/kudig/kudig/pkg/types"

	// Import analyzers to register them
	_ "github.com/kudig/kudig/pkg/analyzer/kernel"
	_ "github.com/kudig/kudig/pkg/analyzer/kubernetes"
	_ "github.com/kudig/kudig/pkg/analyzer/network"
	_ "github.com/kudig/kudig/pkg/analyzer/process"
	_ "github.com/kudig/kudig/pkg/analyzer/runtime"
	_ "github.com/kudig/kudig/pkg/analyzer/system"
)

var (
	version = "2.0.0"
	// Global flags
	verbose    bool
	outputFile string
	format     string

	// Online mode flags
	kubeconfig string
	kubeCtx    string
	nodeName   string
	namespace  string
	allNodes   bool

	// Rules mode flags
	rulesFile string
	rulesDir  string
	listRules bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
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

	// Add deprecated flags for backward compatibility
	rootCmd.Flags().Bool("json", false, "Output JSON format (deprecated, use --format json)")
}

func runOffline(cmd *cobra.Command, args []string) error {
	diagnosePath := args[0]

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
		return fmt.Errorf("offline collector not available")
	}

	// Collect data
	config := collector.NewOfflineConfig(diagnosePath)
	data, err := col.Collect(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to collect data: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "节点信息: %s\n", data.NodeInfo.Hostname)
		fmt.Fprintf(os.Stderr, "\n开始诊断检查...\n\n")
	}

	// Run analyzers
	results, err := analyzer.DefaultRegistry.ExecuteAll(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to run analyzers: %w", err)
	}

	// Collect all issues
	issues := analyzer.CollectIssues(results)

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
	if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
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
		if err := os.WriteFile(outputFile, output, 0644); err != nil {
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
		os.Exit(2)
	} else if len(issues) > 0 {
		os.Exit(1)
	}

	return nil
}

func runLegacy(cmd *cobra.Command, args []string) error {
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
		os.Exit(2)
	} else if report.Summary.Total > 0 {
		os.Exit(1)
	}

	return nil
}

func runListAnalyzers(cmd *cobra.Command, args []string) error {
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

func runOnline(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("online collector not available")
	}

	// Build config
	config := &collector.Config{
		Kubeconfig:     kubeconfig,
		Context:        kubeCtx,
		NodeName:       nodeName,
		Namespace:      namespace,
		AllNodes:       allNodes,
		TimeoutSeconds: 60,
	}

	// Collect data
	if verbose {
		fmt.Fprintf(os.Stderr, "正在连接 Kubernetes 集群...\n")
	}

	data, err := col.Collect(ctx, config)
	if err != nil {
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
		return fmt.Errorf("failed to run analyzers: %w", err)
	}

	// Collect all issues
	issues := analyzer.CollectIssues(results)

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
		if err := os.WriteFile(outputFile, output, 0644); err != nil {
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
		os.Exit(2)
	} else if len(issues) > 0 {
		os.Exit(1)
	}

	return nil
}

func runRules(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("offline collector not available")
	}

	config := collector.NewOfflineConfig(diagnosePath)
	data, err := col.Collect(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to collect data: %w", err)
	}

	// Run rule engine
	engine := rules.NewEngine(loader)
	issues, err := engine.Evaluate(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to evaluate rules: %w", err)
	}

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
		if err := os.WriteFile(outputFile, output, 0644); err != nil {
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
		os.Exit(2)
	} else if len(issues) > 0 {
		os.Exit(1)
	}

	return nil
}
