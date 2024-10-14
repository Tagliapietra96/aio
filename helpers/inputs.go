////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

// File Name: inputs.go
// Created by: Matteo Tagliapietra 2024-10-08
// Last Update: 2024-10-14

// This file contains the input structs and functions to manipulate it.
// Every input ha a runnable function that initializes the input and returns the value.

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

// helpers package contains the helpers used in the application.
package helpers

// imports the necessary packages
// fmt package is used to format strings
// math package is used to perform mathematical operations
// os package is used to interact with the operating system
// strings package is used to manipulate strings
// bubbles package is used to add components to the models
// bubbletea package is used to create terminal applications
// lipgloss package is used to style the terminal output
// log package is used to log messages to the console
import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

//
// Input functions
//

// input struct represents an input field.
type input struct {
	in         textinput.Model
	errStyle   lipgloss.Style
	normStyle  lipgloss.Style
	validation func(string) error
	ph         string
	done       bool
	err        error
}

// Init function initializes the input field.
func (i *input) Init() tea.Cmd {
	return textinput.Blink
}

// Update function updates the input field based on the message received.
func (i *input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		i.err = nil
		i.in.Placeholder = i.ph
		i.in.PlaceholderStyle = i.normStyle

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			i.done = false
			i.in.Blur()
			return i, tea.Quit
		case tea.KeyEnter:
			err := i.validation(i.in.Value())
			if err != nil {
				i.err = err
				i.in.PlaceholderStyle = i.errStyle
				i.in.Placeholder = i.err.Error()
				return i, nil
			} else {
				i.done = true
				i.in.Blur()
				return i, tea.Quit
			}
		}
	}

	i.in, cmd = i.in.Update(msg)

	return i, cmd
}

// View function returns the input field as a string.
func (i *input) View() string {
	return fmt.Sprintf(
		"%s\n\n",
		i.in.View(),
	)
}

// RunInoutWithValidation function initializes the input field and returns the value.
// It takes a string as input placeholder and a validation function and returns a string.
// the validation function takes a string as input and returns an error.
func RunInoutWithValidation(ph string, validation func(string) error) string {
	ti := textinput.New()
	ti.Placeholder = ph
	ti.Focus()
	ti.CharLimit = 156
	norms := ti.PlaceholderStyle
	errs := norms.Foreground(lipgloss.Color("205"))
	m := input{in: ti, normStyle: norms, errStyle: errs, validation: validation, ph: ph, done: false, err: nil}
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err.Error())
	}
	if !m.done {
		os.Exit(0)
	}
	return m.in.Value()
}

// RunInput function initializes the input field and returns the value.
// It takes a string as input placeholder and returns a string.
// the value can't be empty.
func RunInput(ph string) string {
	validation := func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("Please enter a value")
		}
		return nil
	}
	return RunInoutWithValidation(ph, validation)
}

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

//
// Confirm functions
//

// confirm struct represents a confirm field.
type confirm struct {
	in       textinput.Model
	question string
	done     bool
	response bool
}

// Init function initializes the confirm field.
func (c *confirm) Init() tea.Cmd {
	return c.in.Focus()
}

// Update function updates the confirm field based on the message received.
func (c *confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			c.done = false
			c.in.Blur()
			return c, tea.Quit
		}

		switch msg.String() {
		case "y", "Y":
			c.in.SetValue("[Yes!]")
			c.done = true
			c.response = true
			c.in.Blur()
			return c, tea.Quit
		case "n", "N":
			c.in.SetValue("[No..]")
			c.done = true
			c.response = false
			c.in.Blur()
			return c, tea.Quit
		}
	}

	return c, nil
}

// View function returns the confirm field as a string.
func (c *confirm) View() string {
	return fmt.Sprintf(
		"%s %s\n\n",
		c.question,
		c.in.View(),
	)
}

// RunConfirm function initializes the confirm field and returns the value.
// It takes a string as input question and returns a boolean.
func RunConfirm(question string) bool {
	ti := textinput.New()
	ti.Placeholder = "y/n"
	ti.Prompt = ""
	ti.Cursor.SetMode(cursor.CursorHide)
	ti.Focus()
	m := confirm{in: ti, question: question, done: false, response: false}
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err.Error())
	}
	if !m.done {
		os.Exit(0)
	}
	return m.response
}

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

//
// Select functions
//

// selectField struct represents a select field.
type selectField struct {
	index   int
	page    int
	options []string
	done    bool
}

// Init function initializes the select field.
func (s *selectField) Init() tea.Cmd {
	return nil
}

// Update function updates the select field based on the message received.
func (s *selectField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	itemsxpage := 10
	numofpages := int(math.Ceil(float64(len(s.options)) / float64(itemsxpage)))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			s.done = false
			return s, tea.Quit
		case tea.KeyUp:
			if s.index%itemsxpage == 0 {
				if s.page > 1 {
					s.page--
					s.index--
				}
			} else {
				s.index--
			}

			if s.index < 0 {
				s.index = 0
			}
		case tea.KeyDown:
			if s.index%itemsxpage == itemsxpage-1 {
				if s.page < numofpages {
					s.page++
					s.index++
				}
			} else {
				s.index++
			}

			if s.index >= len(s.options) {
				s.index = len(s.options) - 1
			}
		case tea.KeyLeft:
			if s.page > 1 {
				s.page--
				s.index -= itemsxpage
			}

			if s.index < 0 {
				s.index = 0
			}
		case tea.KeyRight:
			if s.page < numofpages {
				s.page++
				s.index += itemsxpage
			}

			if s.index >= len(s.options) {
				s.index = len(s.options) - 1
			}
		case tea.KeyEnter:
			s.done = true
			return s, tea.Quit
		}
	}

	return s, nil
}

// View function returns the select field as a string.
func (s *selectField) View() string {
	view := ""
	itemsxpage := 10
	numofpages := int(math.Ceil(float64(len(s.options)) / float64(itemsxpage)))
	start := (s.page - 1) * itemsxpage
	end := s.page * itemsxpage

	if end > len(s.options) {
		end = len(s.options)
	}

	for i, option := range s.options[start:end] {
		cursor := "  "
		if s.index%itemsxpage == i {
			cursor = "> "
		}
		view += fmt.Sprintf("%s%s\n", cursor, option)
	}

	if numofpages > 1 {
		view += fmt.Sprintf("\nPage %d/%d\n\n", s.page, numofpages)
	}

	return view
}

// RunSelect function initializes the select field and returns the value.
// It takes a slice of strings as input options and returns a string.
func RunSelect(options []string) string {
	m := selectField{index: 0, page: 1, options: options, done: false}
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err.Error())
	}
	if !m.done {
		os.Exit(0)
	}
	return m.options[m.index]
}
