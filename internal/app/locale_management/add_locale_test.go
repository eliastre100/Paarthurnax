package locale_management

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/translation"
	"errors"
	"slices"
	"testing"
)

func TestAddLocale(t *testing.T) {
	t.Run("adds the selected available locale and saves the state", func(t *testing.T) {
		state := &settings.ProjectState{Settings: &settings.Settings{
			SourceLocale:       translation.English,
			DestinationLocales: []translation.Locale{translation.French},
		}}
		repository := &fakeProjectStateRepository{state: state}
		selector := &fakeSelector{selected: translation.German.Name()}

		err := AddLocale(repository, nil, selector)

		if err != nil {
			t.Fatalf("AddLocale() error = %v", err)
		}
		if selector.question != "Select the locale to add as destination" {
			t.Errorf("selection question = %q, want destination locale prompt", selector.question)
		}
		if slices.Contains(selector.choices, translation.English.Name()) || slices.Contains(selector.choices, translation.French.Name()) {
			t.Errorf("selection choices = %v, must exclude source and existing destination locales", selector.choices)
		}
		if !slices.Contains(selector.choices, translation.German.Name()) {
			t.Errorf("selection choices = %v, must include German", selector.choices)
		}
		if len(repository.saved) != 1 || repository.saved[0] != state {
			t.Fatalf("saved states = %#v, want the loaded state saved once", repository.saved)
		}
		if got, want := state.Settings.DestinationLocales, []translation.Locale{translation.French, translation.German}; !slices.Equal(got, want) {
			t.Errorf("destination locales = %v, want %v", got, want)
		}
	})

	t.Run("returns state load error without selecting or saving", func(t *testing.T) {
		loadErr := errors.New("state unavailable")
		repository := &fakeProjectStateRepository{loadErr: loadErr}
		selector := &fakeSelector{}

		err := AddLocale(repository, nil, selector)

		if !errors.Is(err, loadErr) {
			t.Fatalf("AddLocale() error = %v, want %v", err, loadErr)
		}
		if selector.calls != 0 || len(repository.saved) != 0 {
			t.Errorf("selection calls = %d, saved states = %d, want 0 and 0", selector.calls, len(repository.saved))
		}
	})

	t.Run("returns selection error without saving", func(t *testing.T) {
		selectionErr := errors.New("selection cancelled")
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}}
		repository := &fakeProjectStateRepository{state: state}
		selector := &fakeSelector{err: selectionErr}

		err := AddLocale(repository, nil, selector)

		if !errors.Is(err, selectionErr) {
			t.Fatalf("AddLocale() error = %v, want %v", err, selectionErr)
		}
		if len(repository.saved) != 0 || len(state.Settings.DestinationLocales) != 0 {
			t.Errorf("saved states = %d, destination locales = %v, want no changes", len(repository.saved), state.Settings.DestinationLocales)
		}
	})

	t.Run("rejects a locale that is not available without saving", func(t *testing.T) {
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}}
		repository := &fakeProjectStateRepository{state: state}
		selector := &fakeSelector{selected: translation.English.Name()}

		err := AddLocale(repository, nil, selector)

		if err == nil || err.Error() != "English (en) is not a valid choice" {
			t.Fatalf("AddLocale() error = %v, want invalid choice error", err)
		}
		if len(repository.saved) != 0 || len(state.Settings.DestinationLocales) != 0 {
			t.Errorf("saved states = %d, destination locales = %v, want no changes", len(repository.saved), state.Settings.DestinationLocales)
		}
	})

	t.Run("returns save error after appending the locale", func(t *testing.T) {
		saveErr := errors.New("disk full")
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}}
		repository := &fakeProjectStateRepository{state: state, saveErr: saveErr}
		selector := &fakeSelector{selected: translation.German.Name()}

		err := AddLocale(repository, nil, selector)

		if !errors.Is(err, saveErr) {
			t.Fatalf("AddLocale() error = %v, want %v", err, saveErr)
		}
		if len(repository.saved) != 1 || !slices.Equal(state.Settings.DestinationLocales, []translation.Locale{translation.German}) {
			t.Errorf("saved states = %d, destination locales = %v, want one save with German appended", len(repository.saved), state.Settings.DestinationLocales)
		}
	})
}
