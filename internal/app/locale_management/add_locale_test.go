package locale_management

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"
	"errors"
	"reflect"
	"slices"
	"testing"
)

func TestAddLocale(t *testing.T) {
	t.Run("translates source documents into Finnish and saves the state", func(t *testing.T) {
		project := sourceProject(t, "config/locales/en.yml",
			*translation.NewSegment("hello", "Hello!"),
			*translation.NewSegment("test", "Test 2"),
			*translation.NewSegment("test_plural.one", "%{count} one"),
			*translation.NewSegment("test_plural.other", "%{count} other"),
		)
		snapshot, err := sync.NewSnapshotFromProject(project, translation.English)
		if err != nil {
			t.Fatalf("NewSnapshotFromProject() error = %v", err)
		}
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}, Snapshot: snapshot}
		repository := &fakeProjectStateRepository{state: state}
		loader := &fakeProjectLoader{project: project}
		engine := &fakeEngine{result: finnishTranslation}
		reporter := &fakeReporter{}
		store := &fakeDocumentStore{}

		err = AddLocale(repository, loader, engine, reporter, store, &fakeSelector{selected: translation.Finnish.Name()})

		if err != nil {
			t.Fatalf("AddLocale() error = %v", err)
		}
		if loader.calls != 1 {
			t.Errorf("project load calls = %d, want 1", loader.calls)
		}
		if got, want := reporter.changes, []uint{1}; !reflect.DeepEqual(got, want) {
			t.Errorf("reported changes = %v, want %v", got, want)
		}
		if len(engine.calls) != 4 {
			t.Errorf("translation calls = %d, want 4", len(engine.calls))
		}
		if got, want := storedDocumentNames(store), []string{"config/locales/fi.yml"}; !reflect.DeepEqual(got, want) {
			t.Errorf("stored documents = %v, want %v", got, want)
		}
		assertSegment(t, project, "config/locales/fi.yml", "hello", "Hei!")
		assertSegment(t, project, "config/locales/fi.yml", "test", "Testi 2")
		assertSegment(t, project, "config/locales/fi.yml", "test_plural.one", "yksi %{count}")
		assertSegment(t, project, "config/locales/fi.yml", "test_plural.other", "muu %{count}")
		if len(repository.saved) != 1 || repository.saved[0] != state {
			t.Fatalf("saved states = %#v, want the loaded state saved once", repository.saved)
		}
		if got, want := state.Settings.DestinationLocales, []translation.Locale{translation.Finnish}; !slices.Equal(got, want) {
			t.Errorf("destination locales = %v, want %v", got, want)
		}
	})

	t.Run("adds the selected available locale and saves an unchanged project", func(t *testing.T) {
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English, DestinationLocales: []translation.Locale{translation.French}}, Snapshot: sync.NewSnapshot("en")}
		repository := &fakeProjectStateRepository{state: state}
		selector := &fakeSelector{selected: translation.German.Name()}

		err := AddLocale(repository, &fakeProjectLoader{project: translation.NewProject()}, &fakeEngine{result: finnishTranslation}, &fakeReporter{}, &fakeDocumentStore{}, selector)

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
	})

	t.Run("returns state load error without loading, selecting, or saving", func(t *testing.T) {
		loadErr := errors.New("state unavailable")
		repository := &fakeProjectStateRepository{loadErr: loadErr}
		loader := &fakeProjectLoader{}
		selector := &fakeSelector{}

		err := AddLocale(repository, loader, &fakeEngine{result: finnishTranslation}, &fakeReporter{}, &fakeDocumentStore{}, selector)

		if !errors.Is(err, loadErr) {
			t.Fatalf("AddLocale() error = %v, want %v", err, loadErr)
		}
		if loader.calls != 0 || selector.calls != 0 || len(repository.saved) != 0 {
			t.Errorf("loader calls = %d, selection calls = %d, saved states = %d, want 0, 0, 0", loader.calls, selector.calls, len(repository.saved))
		}
	})

	t.Run("returns project load error without selecting or saving", func(t *testing.T) {
		loadErr := errors.New("project unavailable")
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}}
		repository := &fakeProjectStateRepository{state: state}
		selector := &fakeSelector{}

		err := AddLocale(repository, &fakeProjectLoader{err: loadErr}, &fakeEngine{result: finnishTranslation}, &fakeReporter{}, &fakeDocumentStore{}, selector)

		if !errors.Is(err, loadErr) {
			t.Fatalf("AddLocale() error = %v, want %v", err, loadErr)
		}
		if selector.calls != 0 || len(repository.saved) != 0 {
			t.Errorf("selection calls = %d, saved states = %d, want 0 and 0", selector.calls, len(repository.saved))
		}
	})

	t.Run("returns selection errors without saving", func(t *testing.T) {
		selectionErr := errors.New("selection cancelled")
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}}
		repository := &fakeProjectStateRepository{state: state}

		err := AddLocale(repository, &fakeProjectLoader{project: translation.NewProject()}, &fakeEngine{result: finnishTranslation}, &fakeReporter{}, &fakeDocumentStore{}, &fakeSelector{err: selectionErr})

		if !errors.Is(err, selectionErr) {
			t.Fatalf("AddLocale() error = %v, want %v", err, selectionErr)
		}
		if len(repository.saved) != 0 || len(state.Settings.DestinationLocales) != 0 {
			t.Errorf("saved states = %d, destination locales = %v, want no changes", len(repository.saved), state.Settings.DestinationLocales)
		}
	})

	t.Run("rejects unavailable locales without saving", func(t *testing.T) {
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}}
		repository := &fakeProjectStateRepository{state: state}

		err := AddLocale(repository, &fakeProjectLoader{project: translation.NewProject()}, &fakeEngine{result: finnishTranslation}, &fakeReporter{}, &fakeDocumentStore{}, &fakeSelector{selected: translation.English.Name()})

		if err == nil || err.Error() != "English (en) is not a valid choice" {
			t.Fatalf("AddLocale() error = %v, want invalid choice error", err)
		}
		if len(repository.saved) != 0 || len(state.Settings.DestinationLocales) != 0 {
			t.Errorf("saved states = %d, destination locales = %v, want no changes", len(repository.saved), state.Settings.DestinationLocales)
		}
	})

	t.Run("does not save the locale when translation fails", func(t *testing.T) {
		project := sourceProject(t, "config/locales/en.yml", *translation.NewSegment("hello", "Hello!"))
		snapshot, err := sync.NewSnapshotFromProject(project, translation.English)
		if err != nil {
			t.Fatalf("NewSnapshotFromProject() error = %v", err)
		}
		translateErr := errors.New("translator unavailable")
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}, Snapshot: snapshot}
		repository := &fakeProjectStateRepository{state: state}

		err = AddLocale(repository, &fakeProjectLoader{project: project}, &fakeEngine{result: func(translation.Segment, translation.Locale) (translation.Segment, error) {
			return translation.Segment{}, translateErr
		}}, &fakeReporter{}, &fakeDocumentStore{}, &fakeSelector{selected: translation.Finnish.Name()})

		if !errors.Is(err, translateErr) {
			t.Fatalf("AddLocale() error = %v, want %v", err, translateErr)
		}
		if len(repository.saved) != 0 {
			t.Errorf("saved states = %d, want 0", len(repository.saved))
		}
	})

	t.Run("returns save errors after successful translation", func(t *testing.T) {
		saveErr := errors.New("disk full")
		state := &settings.ProjectState{Settings: &settings.Settings{SourceLocale: translation.English}, Snapshot: sync.NewSnapshot("en")}
		repository := &fakeProjectStateRepository{state: state, saveErr: saveErr}

		err := AddLocale(repository, &fakeProjectLoader{project: translation.NewProject()}, &fakeEngine{result: finnishTranslation}, &fakeReporter{}, &fakeDocumentStore{}, &fakeSelector{selected: translation.German.Name()})

		if !errors.Is(err, saveErr) {
			t.Fatalf("AddLocale() error = %v, want %v", err, saveErr)
		}
		if len(repository.saved) != 1 || !slices.Equal(state.Settings.DestinationLocales, []translation.Locale{translation.German}) {
			t.Errorf("saved states = %d, destination locales = %v, want one save with German appended", len(repository.saved), state.Settings.DestinationLocales)
		}
	})
}

func sourceProject(t *testing.T, name string, segments ...translation.Segment) *translation.Project {
	t.Helper()
	document := translation.NewDocument(name)
	for _, segment := range segments {
		if err := document.SetSegment(translation.English, segment); err != nil {
			t.Fatalf("SetSegment() error = %v", err)
		}
	}
	project := translation.NewProject()
	project.AddDocument(document)
	return project
}

func finnishTranslation(segment translation.Segment, _ translation.Locale) (translation.Segment, error) {
	translations := map[string]string{
		"hello":             "Hei!",
		"test":              "Testi 2",
		"test_plural.one":   "yksi %{count}",
		"test_plural.other": "muu %{count}",
	}
	return *translation.NewSegment(segment.Key, translations[segment.Key]), nil
}

func assertSegment(t *testing.T, project *translation.Project, documentName, key, want string) {
	t.Helper()
	document, err := project.GetDocumentIn(translation.Finnish, documentName)
	if err != nil {
		t.Fatalf("GetDocumentIn() error = %v", err)
	}
	segment, err := document.GetSegment(translation.Finnish, key)
	if err != nil {
		t.Fatalf("GetSegment(%q) error = %v", key, err)
	}
	if segment.Value != want {
		t.Errorf("segment %q = %q, want %q", key, segment.Value, want)
	}
}

func storedDocumentNames(store *fakeDocumentStore) []string {
	names := make([]string, 0, len(store.stored))
	for _, document := range store.stored {
		names = append(names, document.Name)
	}
	return names
}
