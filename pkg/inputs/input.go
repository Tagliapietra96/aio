// inputs input functions
package inputs

import (
	"aio/pkg/log"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
				i.in.SetValue("")
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

// RunInputWithValidation function initializes the input field and returns the value.
// It takes a string as input placeholder and a validation function and returns a string.
// the validation function takes a string as input and returns an error.
func RunInputWithValidation(ph string, validation func(string) error) string {
	ti := textinput.New()
	ti.Placeholder = ph
	ti.Focus()
	ti.CharLimit = 156
	norms := ti.PlaceholderStyle
	errs := norms.Foreground(lipgloss.Color("196"))
	m := input{in: ti, normStyle: norms, errStyle: errs, validation: validation, ph: ph, done: false, err: nil}
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		log.Fat(err)
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
			return errors.New("please enter a value")
		}
		return nil
	}
	return RunInputWithValidation(ph, validation)
}
