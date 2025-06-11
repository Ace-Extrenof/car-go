package main

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

// DashboardResult holds the settings user chose
type DashboardResult struct {
	Randomised bool
	NumOrders  int
}

type model struct {
	cursor     int
	choices    []string
	randomised bool

	// input mode
	typingNum bool
	textInput textinput.Model

	numOrders string
	finished  bool
	err       error
}

func ShowDashboard() (*DashboardResult, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.StartReturningModel()
	if err != nil {
		return nil, err
	}
	mod := m.(model)

	if mod.finished {
		num, err := strconv.Atoi(mod.numOrders)
		if err != nil || num < 1 || num > 100 {
			return nil, fmt.Errorf("invalid number of orders: %s", mod.numOrders)
		}
		return &DashboardResult{
			Randomised: mod.randomised,
			NumOrders:  num,
		}, nil
	}
	return nil, nil // user quit
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "1-100"
	ti.CharLimit = 3
	ti.Width = 10

	return model{
		cursor:     0,
		choices:    []string{"Randomised Orders: OFF", "Number of Orders: 1", "Confirm & Start", "Exit"},
		randomised: false,
		numOrders:  "1",
		typingNum:  false,
		textInput:  ti,
		finished:   false,
		err:        nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.typingNum {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				num, err := strconv.Atoi(m.textInput.Value())
				if err != nil || num < 1 || num > 100 {
					m.err = fmt.Errorf("please enter a number between 1 and 100")
					return m, cmd
				}

				m.numOrders = m.textInput.Value()
				m.choices[1] = "Number of Orders: " + m.numOrders
				m.err = nil
				m.typingNum = false
				return m, nil

			case "esc":
				m.typingNum = false
				m.err = nil
				return m, nil
			}
		}

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			switch m.cursor {
			case 0:
				m.randomised = !m.randomised
				if m.randomised {
					m.choices[0] = "Randomised Orders: ON"
				} else {
					m.choices[0] = "Randomised Orders: OFF"
				}
			case 1:
				m.typingNum = true
				m.textInput.SetValue(m.numOrders)
				m.textInput.Focus()
			case 2:
				m.finished = true
				return m, tea.Quit
			case 3:
				m.finished = false
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.typingNum {
		return fmt.Sprintf(
			"\n📦 Enter number of orders (1-100):\n\n%s\n\n(Enter to confirm, Esc to cancel)\n%s\n",
			m.textInput.View(),
			func() string {
				if m.err != nil {
					return fmt.Sprintf("\n❌ %s\n", m.err)
				}
				return ""
			}(),
		)
	}

	s := "\n📦 Dashboard - Configure your session\n\n"

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	if m.err != nil {
		s += fmt.Sprintf("\n❌ %s\n", m.err)
	}

	s += "\nUse ↑/↓ to navigate, Enter to select.\nPress q or Ctrl+C to quit.\n"

	return s
}

