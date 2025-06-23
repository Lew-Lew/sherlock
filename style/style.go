package style

import "github.com/charmbracelet/lipgloss"

var IssueKey = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("15"))

var IssueAlert = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("1"))

var IssueInfo = lipgloss.NewStyle().
	Foreground(lipgloss.Color("7"))
