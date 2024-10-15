// style package is used to style the terminal output
package style

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Print render the string in a new line.
func Print(s string) {
	fmt.Println(s)
}

// Title render the title style.
func Title(s string) string {
	t := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15"))
	return t.Render(s)
}

// PrintTitle render the title style.
func PrintTitle(s string) {
	t := Title(s)
	Print(t)
}
