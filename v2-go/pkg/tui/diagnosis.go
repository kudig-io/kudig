package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/kudig/kudig/pkg/types"
)

// DiagnosisCompleteMsg is sent when diagnosis is complete
type DiagnosisCompleteMsg struct {
	Issues []types.Issue
	Data   *types.DiagnosticData
	Error  error
}

// startDiagnosis starts the diagnosis process
func (m Model) startDiagnosis() tea.Cmd {
	return func() tea.Msg {
		// This is a placeholder that returns an empty result
		// In a real implementation, this would:
		// 1. Get the appropriate collector based on mode
		// 2. Collect diagnostic data
		// 3. Run analyzers
		// 4. Return the results

		// For now, return empty result to allow compilation
		return DiagnosisCompleteMsg{
			Issues: []types.Issue{},
			Data:   nil,
			Error:  nil,
		}
	}
}

// RunTUI starts the terminal user interface
func RunTUI() error {
	model := NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
