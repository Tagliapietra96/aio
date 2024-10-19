// inputs select functions
package inputs

import (
	"aio/pkg/log"
	"fmt"
	"math"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

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
		log.Fat(err)
	}
	if !m.done {
		os.Exit(0)
	}
	return m.options[m.index]
}
