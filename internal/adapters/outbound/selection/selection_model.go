package selection

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

type selectionModel struct {
	question string
	choices  []string
	query    string
	cursor   int
	height   int
	selected bool
	aborted  bool
	complete bool
}

func newSelectionModel(question string, choices []string) selectionModel {
	return selectionModel{
		question: question,
		choices:  choices,
		height:   defaultVisibleChoices + 5,
	}
}

func (m selectionModel) Init() tea.Cmd {
	return nil
}

func (m selectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.height = size.Height
		return m, nil
	}

	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	if key.Mod == tea.ModCtrl && key.Code == 'c' || key.Code == tea.KeyEsc {
		m.aborted = true
		m.complete = true
		return m, tea.Quit
	}

	return m.handleKey(key)
}

func (m selectionModel) handleKey(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	visible := m.filtered()

	switch key.Code {
	case tea.KeyEnter:
		if len(visible) > 0 {
			m.selected = true
			m.complete = true
			return m, tea.Quit
		}
	case tea.KeyUp:
		if len(visible) > 0 {
			m.cursor = max(m.cursor-1, 0)
		}
	case tea.KeyDown:
		if len(visible) > 0 {
			m.cursor = min(m.cursor+1, len(visible)-1)
		}
	case tea.KeyBackspace:
		m.removeLastQueryRune()
	default:
		if key.Text != "" {
			m.query += key.Text
			m.cursor = 0
		}
	}
	return m, nil
}

func (m *selectionModel) removeLastQueryRune() {
	query := []rune(m.query)
	if len(query) == 0 {
		return
	}

	m.query = string(query[:len(query)-1])
	m.cursor = 0
}

func (m selectionModel) result() (string, error) {
	if m.aborted || !m.selected {
		return "", ErrAborted
	}
	return m.filtered()[m.cursor], nil
}

func (m selectionModel) filtered() []string {
	if m.query == "" {
		return m.choices
	}

	query := strings.ToLower(m.query)
	filtered := make([]string, 0, len(m.choices))
	for _, choice := range m.choices {
		if strings.Contains(strings.ToLower(choice), query) {
			filtered = append(filtered, choice)
		}
	}
	return filtered
}
