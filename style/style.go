///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// File Name: style.go
// Created by: Matteo Tagliapietra 2024-10-15
// Last Update: 2024-10-15

// This file contains the style used in the application.
// It contains all the function to display the output in a fancy way.

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

// style package is used to style the terminal output
package style

// imports the necessary packages
// fmt package is used to format strings
// lipgloss package is used to style the terminal output
import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

//
// Style functions
//

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
