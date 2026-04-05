package analyzer

import "github.com/kudig/kudig/pkg/analyzer"

func init() {
	// Register eBPF analyzers
	_ = analyzer.Register(NewTCPAnalyzer())
	_ = analyzer.Register(NewDNSAnalyzer())
	_ = analyzer.Register(NewFileIOAnalyzer())
}
