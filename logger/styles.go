// logger package styles
package logger

import "github.com/charmbracelet/lipgloss"

// colors
var (
	BrigthColor = lipgloss.Color("15")
	ErrorColor  = lipgloss.Color("196")
)

// styles
var (
	TitleStyle = lipgloss.NewStyle().Foreground(BrigthColor).Bold(true)
	ErrorStyle = lipgloss.NewStyle().Foreground(ErrorColor).Bold(true)
)
