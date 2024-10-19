// inputs confirm functions
package inputs

import (
	"aio/pkg/log"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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
		log.Fat(err)
	}
	if !m.done {
		os.Exit(0)
	}
	return m.response
}
