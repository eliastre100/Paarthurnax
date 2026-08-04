package selection

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

const selectorChromeHeight = 5

func (m selectionModel) View() tea.View {
	if m.complete {
		return tea.NewView("")
	}

	var view strings.Builder
	fmt.Fprintf(&view, "%s\nSearch: %s\n\n", m.question, m.query)
	m.writeChoices(&view)
	m.padView(&view)
	view.WriteString("[up/down navigate] [enter choose] [esc cancel]")
	return tea.NewView(view.String())
}

func (m selectionModel) writeChoices(view *strings.Builder) {
	filtered := m.filtered()
	start, end := m.visibleRange(len(filtered))
	for index := start; index < end; index++ {
		pointer := " "
		if index == m.cursor {
			pointer = ">"
		}
		fmt.Fprintf(view, "%s %s\n", pointer, filtered[index])
	}

	if len(filtered) == 0 {
		view.WriteString("  No matching choices\n")
	} else if hidden := len(filtered) - (end - start); hidden > 0 {
		fmt.Fprintf(view, "  (%d more results)\n", hidden)
	}
}

func (m selectionModel) padView(view *strings.Builder) {
	for lineCount(view.String()) < m.height-1 {
		view.WriteByte('\n')
	}
}

func (m selectionModel) visibleRange(total int) (int, int) {
	capacity := max(m.height-selectorChromeHeight, 1)
	if total <= capacity {
		return 0, total
	}

	start := min(max(m.cursor-capacity/2, 0), total-capacity)
	return start, start + capacity
}

func lineCount(text string) int {
	return strings.Count(text, "\n") + 1
}
