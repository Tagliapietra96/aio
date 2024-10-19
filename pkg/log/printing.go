// log package printing functions
package log

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Print function prints a message to the console
// with a unformatted style
func Print(msg string, args ...any) {
	style := lipgloss.NewStyle()
	PrintS(msg, style, args...)
}

// PrintS function prints a message to the console
// with the specified style
func PrintS(msg string, style lipgloss.Style, args ...any) {
	s := fmt.Sprintf(msg, args...)
	s = style.Render(s)
	println(s)
}
