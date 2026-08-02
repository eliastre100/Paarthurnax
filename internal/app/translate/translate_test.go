package translate

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"
	"errors"
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("synchronizes new source documents, reports progress, and saves the resulting snapshot", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json",
			translation.Segment{Key: "farewell", Value: "Goodbye"},
			translation.Segment{Key: "greeting", Value: "Hello"},
		)
		loader := &fakeProjectLoader{project: project}
		repository := &fakeProjectStateRepository{state: &settings.ProjectState{Settings: testSettings(), Snapshot: sync.NewSnapshot("en")}}
		store := &fakeDocumentStore{}
		engine := &fakeEngine{result: translatedSegment}
		reporter := &fakeReporter{}

		err := Execute(loader, store, repository, engine, reporter)

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if loader.calls != 1 {
			t.Errorf("project load calls = %d, want 1", loader.calls)
		}
		if !reflect.DeepEqual(reporter.changes, []uint{1}) || !reflect.DeepEqual(reporter.documents, []string{source.Name}) {
			t.Errorf("reported changes = %#v, documents = %#v, want [1], [%q]", reporter.changes, reporter.documents, source.Name)
		}
		if got, want := len(engine.calls), 4; got != want {
			t.Errorf("translation calls = %d, want %d", got, want)
		}
		if got, want := documentNames(store.stored), []string{"locales/fr/messages.json", "locales/ja/messages.json"}; !reflect.DeepEqual(got, want) {
			t.Errorf("stored documents = %v, want %v", got, want)
		}
		assertSegment(t, project, "locales/fr/messages.json", translation.French, "greeting", "Hello in fr")
		assertSegment(t, project, "locales/ja/messages.json", translation.Japanese, "farewell", "Goodbye in ja")
		if len(repository.saved) != 1 || repository.saved[0] != repository.state {
			t.Fatalf("saved states = %#v, want the loaded state once", repository.saved)
		}
		if got := repository.state.Snapshot; got == nil || got.Locale != "en" || len(got.Documents) != 1 {
			t.Errorf("saved snapshot = %#v, want one English source document", got)
		}
	})

	t.Run("records a fresh snapshot even when source files are unchanged", func(t *testing.T) {
		project, _ := projectWithSource(t, "locales/en/messages.json", translation.Segment{Key: "greeting", Value: "Hello"})
		snapshot, err := sync.NewSnapshotFromProject(project, translation.English)
		if err != nil {
			t.Fatalf("NewSnapshotFromProject() error = %v", err)
		}
		repository := &fakeProjectStateRepository{state: &settings.ProjectState{Settings: testSettings(), Snapshot: snapshot}}
		store := &fakeDocumentStore{}
		engine := &fakeEngine{result: translatedSegment}
		reporter := &fakeReporter{}

		err = Execute(&fakeProjectLoader{project: project}, store, repository, engine, reporter)

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if !reflect.DeepEqual(reporter.changes, []uint{0}) || len(reporter.documents) != 0 {
			t.Errorf("reported changes = %#v, documents = %#v, want [0], []", reporter.changes, reporter.documents)
		}
		if len(engine.calls) != 0 || len(store.stored) != 0 || len(repository.saved) != 1 {
			t.Errorf("side effects: translations=%d stores=%d saved=%d, want 0, 0, 1", len(engine.calls), len(store.stored), len(repository.saved))
		}
	})

	t.Run("does not save state when processing a change fails", func(t *testing.T) {
		project, _ := projectWithSource(t, "locales/en/messages.json", translation.Segment{Key: "greeting", Value: "Hello"})
		translateErr := errors.New("translator unavailable")
		repository := &fakeProjectStateRepository{state: &settings.ProjectState{Settings: testSettings(), Snapshot: sync.NewSnapshot("en")}}
		reporter := &fakeReporter{}

		err := Execute(&fakeProjectLoader{project: project}, &fakeDocumentStore{}, repository, &fakeEngine{result: func(translation.Segment, translation.Locale) (translation.Segment, error) {
			return translation.Segment{}, translateErr
		}}, reporter)

		if !errors.Is(err, translateErr) {
			t.Fatalf("Execute() error = %v, want %v", err, translateErr)
		}
		if len(repository.saved) != 0 {
			t.Errorf("saved states = %d, want 0", len(repository.saved))
		}
	})
}
