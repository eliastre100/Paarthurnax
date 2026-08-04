package project_initialization

import (
	"Paarthurnax/internal/domain/translation"
	"errors"
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("saves source locale and snapshot", func(t *testing.T) {
		project := projectWithCatalogs(t,
			map[translation.Locale][]translation.Segment{
				translation.English: {{Key: "greeting", Value: "Hello"}},
				translation.French:  {{Key: "greeting", Value: "Bonjour"}},
			},
			map[translation.Locale][]translation.Segment{
				translation.English: {{Key: "farewell", Value: "Goodbye"}},
			},
		)
		repository := &fakeProjectStateRepository{}
		selector := &fakeSelector{selected: "en"}

		err := Execute(&fakeProjectLoader{project: project}, repository, selector)

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if selector.question != "Select the source locale" {
			t.Errorf("selection question = %q, want Select the source locale", selector.question)
		}
		if !sameStrings(selector.choices, []string{"en", "fr"}) {
			t.Errorf("selection choices = %v, want locales en and fr", selector.choices)
		}
		if len(repository.saved) != 1 {
			t.Fatalf("saved states = %d, want 1", len(repository.saved))
		}

		state := repository.saved[0]
		if state.Settings.SourceLocale != translation.English || len(state.Settings.DestinationLocales) != 0 {
			t.Errorf("settings = %#v, want English source and no destinations", state.Settings)
		}
		if state.Snapshot.Locale != "en" || len(state.Snapshot.Documents) != 2 {
			t.Errorf("snapshot = %#v, want English snapshot with two documents", state.Snapshot)
		}
	})

	t.Run("returns loader error without selecting or saving", func(t *testing.T) {
		loadErr := errors.New("project unavailable")
		repository := &fakeProjectStateRepository{}
		selector := &fakeSelector{}

		err := Execute(&fakeProjectLoader{err: loadErr}, repository, selector)

		if !errors.Is(err, loadErr) {
			t.Fatalf("Execute() error = %v, want %v", err, loadErr)
		}
		if selector.calls != 0 || len(repository.saved) != 0 {
			t.Errorf("selection calls = %d, saved states = %d, want 0 and 0", selector.calls, len(repository.saved))
		}
	})

	t.Run("selects the only project locale without prompting", func(t *testing.T) {
		repository := &fakeProjectStateRepository{}
		selector := &fakeSelector{}
		project := projectWithCatalogs(t, map[translation.Locale][]translation.Segment{
			translation.English: nil,
		})

		err := Execute(&fakeProjectLoader{project: project}, repository, selector)

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if selector.calls != 0 {
			t.Errorf("selection calls = %d, want 0", selector.calls)
		}
		if len(repository.saved) != 1 {
			t.Fatalf("saved states = %d, want 1", len(repository.saved))
		}
		if got := repository.saved[0].Settings.SourceLocale; got != translation.English {
			t.Errorf("source locale = %q, want %q", got, translation.English)
		}
	})

	t.Run("returns selection error without saving", func(t *testing.T) {
		selectionErr := errors.New("selection cancelled")
		repository := &fakeProjectStateRepository{}
		selector := &fakeSelector{err: selectionErr}
		project := projectWithCatalogs(t, map[translation.Locale][]translation.Segment{
			translation.English: nil,
			translation.French:  nil,
		})

		err := Execute(&fakeProjectLoader{project: project}, repository, selector)

		if !errors.Is(err, selectionErr) {
			t.Fatalf("Execute() error = %v, want %v", err, selectionErr)
		}
		if len(repository.saved) != 0 {
			t.Errorf("saved states = %d, want 0", len(repository.saved))
		}
	})

	t.Run("rejects a locale outside the project", func(t *testing.T) {
		repository := &fakeProjectStateRepository{}
		selector := &fakeSelector{selected: "fr"}
		project := projectWithCatalogs(t, map[translation.Locale][]translation.Segment{
			translation.English: nil,
			translation.German:  nil,
		})

		err := Execute(&fakeProjectLoader{project: project}, repository, selector)

		if err == nil || err.Error() != "fr is not a valid locale" {
			t.Fatalf("Execute() error = %v, want invalid locale error", err)
		}
		if len(repository.saved) != 0 {
			t.Errorf("saved states = %d, want 0", len(repository.saved))
		}
	})

	t.Run("returns save error after building state", func(t *testing.T) {
		saveErr := errors.New("disk full")
		repository := &fakeProjectStateRepository{saveErr: saveErr}
		selector := &fakeSelector{selected: "en"}

		err := Execute(&fakeProjectLoader{project: projectWithCatalogs(t, map[translation.Locale][]translation.Segment{translation.English: nil})}, repository, selector)

		if !errors.Is(err, saveErr) {
			t.Fatalf("Execute() error = %v, want %v", err, saveErr)
		}
		if len(repository.saved) != 1 || repository.saved[0].Settings.SourceLocale != translation.English {
			t.Errorf("saved states = %#v, want one English state", repository.saved)
		}
	})
}

func projectWithCatalogs(t *testing.T, catalogs ...map[translation.Locale][]translation.Segment) *translation.Project {
	t.Helper()
	project := translation.NewProject()
	for index, documentCatalogs := range catalogs {
		document := translation.NewDocument(string(rune('a'+index)) + ".json")
		for locale, segments := range documentCatalogs {
			catalog := translation.NewCatalog(locale)
			for _, segment := range segments {
				catalog.AddSegment(segment)
			}
			if err := document.AddCatalog(catalog); err != nil {
				t.Fatalf("AddCatalog() error = %v", err)
			}
		}
		project.AddDocument(document)
	}
	return project
}

func sameStrings(got, want []string) bool {
	return reflect.DeepEqual(mapStrings(got), mapStrings(want))
}

func mapStrings(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}
