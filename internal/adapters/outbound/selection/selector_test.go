package selection

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestExecuteRejectsEmptyChoices(t *testing.T) {
	_, err := execute("Colour", nil, strings.NewReader(""), &bytes.Buffer{}, false)
	if !errors.Is(err, ErrNoChoices) {
		t.Fatalf("execute() error = %v, want %v", err, ErrNoChoices)
	}
}

func TestSelectorSelectsChoice(t *testing.T) {
	var output bytes.Buffer
	selector := newSelector(strings.NewReader("green\n"), &output, false)

	choice, err := selector.Select("Colour", []string{"red", "green"})
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if choice != "green" {
		t.Errorf("Select() = %q, want green", choice)
	}
}

func TestPromptValidatesThenReturnsChoice(t *testing.T) {
	var output bytes.Buffer
	choice, err := execute("Colour", []string{"red", "green"}, strings.NewReader("blue\n green \n"), &output, false)
	if err != nil {
		t.Fatalf("execute() error = %v", err)
	}
	if choice != "green" {
		t.Errorf("choice = %q, want green", choice)
	}
	if got := output.String(); !strings.Contains(got, "Invalid selection") {
		t.Errorf("output = %q, want validation message", got)
	}
}

func TestPromptReturnsAbortedAtEndOfInput(t *testing.T) {
	_, err := execute("Colour", []string{"red"}, strings.NewReader(""), &bytes.Buffer{}, false)
	if !errors.Is(err, ErrAborted) {
		t.Fatalf("execute() error = %v, want %v", err, ErrAborted)
	}
}

func TestSelectorFiltersNavigatesAndSelects(t *testing.T) {
	model := newSelectionModel("Choice", []string{"alpha", "alpine", "banana"})
	model = updateSelector(t, model, tea.Key{Text: "a", Code: 'a'})
	model = updateSelector(t, model, tea.Key{Text: "l", Code: 'l'})
	if got := model.filtered(); strings.Join(got, ",") != "alpha,alpine" {
		t.Fatalf("filtered choices = %v, want [alpha alpine]", got)
	}

	model = updateSelector(t, model, tea.Key{Code: tea.KeyDown})
	if !strings.Contains(model.View().Content, "> alpine") {
		t.Errorf("view does not point to selected choice:\n%s", model.View().Content)
	}
	model = updateSelector(t, model, tea.Key{Code: tea.KeyEnter})
	if !model.selected || model.filtered()[model.cursor] != "alpine" {
		t.Errorf("selection = selected:%t choice:%q, want selected: true choice: alpine", model.selected, model.filtered()[model.cursor])
	}
	if got := model.View().Content; got != "" {
		t.Errorf("completed view = %q, want empty", got)
	}
}

func TestSelectorCanBeCancelled(t *testing.T) {
	model := updateSelector(t, newSelectionModel("Choice", []string{"alpha"}), tea.Key{Code: tea.KeyEsc})
	if !model.aborted {
		t.Error("aborted = false, want true")
	}
	if got := model.View().Content; got != "" {
		t.Errorf("cancelled view = %q, want empty", got)
	}
}

func TestWriteResult(t *testing.T) {
	var output bytes.Buffer
	if err := writeResult(&output, "Default locale", "<canceled>"); err != nil {
		t.Fatalf("writeResult() error = %v", err)
	}
	if got := output.String(); got != "Default locale: <canceled>\n" {
		t.Errorf("output = %q, want cancellation result", got)
	}
}

func TestWriteCancellationResult(t *testing.T) {
	var output bytes.Buffer
	err := writeCancellationResult(&output, "Default locale", ErrAborted)
	if !errors.Is(err, ErrAborted) {
		t.Fatalf("writeCancellationResult() error = %v, want %v", err, ErrAborted)
	}
	if got := output.String(); got != "Default locale: <canceled>\n" {
		t.Errorf("output = %q, want cancellation result", got)
	}
}

func TestWriteCancellationResultReturnsOriginalError(t *testing.T) {
	original := errors.New("selection failed")
	err := writeCancellationResult(&bytes.Buffer{}, "Default locale", original)
	if !errors.Is(err, original) {
		t.Errorf("writeCancellationResult() error = %v, want %v", err, original)
	}
}

func TestSelectorCapsResultsAndShowsHiddenCount(t *testing.T) {
	model := newSelectionModel("Choice", []string{"one", "two", "three", "four", "five"})
	updated, _ := model.Update(tea.WindowSizeMsg{Height: 8})
	model = updated.(selectionModel)

	view := model.View().Content
	if !strings.Contains(view, "> one\n  two\n  three\n") {
		t.Errorf("view does not contain the capped initial choices:\n%s", view)
	}
	if !strings.Contains(view, "(2 more results)") {
		t.Errorf("view does not report hidden choices:\n%s", view)
	}
	if strings.Contains(view, "four\n") {
		t.Errorf("view contains choices beyond the result cap:\n%s", view)
	}
	lines := strings.Split(view, "\n")
	if got := lines[len(lines)-1]; got != "[up/down navigate] [enter choose] [esc cancel]" {
		t.Errorf("bottom line = %q, want selection help", got)
	}

	for range 4 {
		model = updateSelector(t, model, tea.Key{Code: tea.KeyDown})
	}
	view = model.View().Content
	if !strings.Contains(view, "> five") {
		t.Errorf("view does not scroll to the cursor:\n%s", view)
	}
	if !strings.Contains(view, "(2 more results)") {
		t.Errorf("view does not retain hidden result count after scrolling:\n%s", view)
	}
}

func updateSelector(t *testing.T, model selectionModel, key tea.Key) selectionModel {
	t.Helper()
	updated, _ := model.Update(tea.KeyPressMsg(key))
	return updated.(selectionModel)
}
