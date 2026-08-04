package selection

import (
	"strings"
	"testing"
)

func TestSelectorViewShowsEmptyStateAndReservesTerminalBottomLine(t *testing.T) {
	model := newSelectionModel("Choice", []string{"alpha"})
	model.query = "z"
	model.height = 8

	view := model.View().Content
	if !strings.Contains(view, "Choice\nSearch: z\n\n  No matching choices\n") {
		t.Errorf("view does not show the empty state:\n%s", view)
	}
	lines := strings.Split(view, "\n")
	if got, want := len(lines), model.height-1; got != want {
		t.Errorf("view line count = %d, want %d", got, want)
	}
	if got := lines[len(lines)-1]; got != "[up/down navigate] [enter choose] [esc cancel]" {
		t.Errorf("bottom line = %q, want selection help", got)
	}
}

func TestSelectorVisibleRangeKeepsCursorInView(t *testing.T) {
	model := selectionModel{height: 8, cursor: 4}
	start, end := model.visibleRange(5)
	if start != 2 || end != 5 {
		t.Errorf("visibleRange() = (%d, %d), want (2, 5)", start, end)
	}
}
