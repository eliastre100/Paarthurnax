package translate

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestWorkerUpdateSegment(t *testing.T) {
	t.Run("translates every destination and creates missing documents", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json", translation.Segment{Key: "greeting", Value: "Hello"})
		engine := &fakeEngine{result: translatedSegment}
		worker := NewWorker(testSettings(), project, engine, &fakeDocumentStore{}, &fakeReporter{})

		if err := worker.UpdateSegment(source.Name, "greeting"); err != nil {
			t.Fatalf("UpdateSegment() error = %v", err)
		}

		wantCalls := []translationCall{
			{segment: source.Catalog(translation.English).Segments["greeting"], source: translation.English, target: translation.French},
			{segment: source.Catalog(translation.English).Segments["greeting"], source: translation.English, target: translation.Japanese},
		}
		if !reflect.DeepEqual(engine.calls, wantCalls) {
			t.Errorf("engine calls = %#v, want %#v", engine.calls, wantCalls)
		}
		assertSegment(t, project, "locales/fr/messages.json", translation.French, "greeting", "Hello in fr")
		assertSegment(t, project, "locales/ja/messages.json", translation.Japanese, "greeting", "Hello in ja")
	})

	t.Run("updates an existing destination document", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json", translation.Segment{Key: "greeting", Value: "Hello"})
		existing := translation.NewDocument("locales/fr/messages.json")
		mustSetSegment(t, existing, translation.French, translation.Segment{Key: "greeting", Value: "old"})
		project.AddDocument(existing)
		settings := testSettings()
		settings.DestinationLocales = []translation.Locale{translation.French}
		worker := NewWorker(settings, project, &fakeEngine{result: translatedSegment}, &fakeDocumentStore{}, &fakeReporter{})

		if err := worker.UpdateSegment(source.Name, "greeting"); err != nil {
			t.Fatalf("UpdateSegment() error = %v", err)
		}

		if project.Documents[existing.Name] != existing {
			t.Error("UpdateSegment() replaced the existing destination document")
		}
		assertSegment(t, project, "locales/fr/messages.json", translation.French, "greeting", "Hello in fr")
	})

	t.Run("returns contextual errors without translating unavailable source data", func(t *testing.T) {
		engine := &fakeEngine{result: translatedSegment}
		worker := NewWorker(testSettings(), translation.NewProject(), engine, &fakeDocumentStore{}, &fakeReporter{})

		err := worker.UpdateSegment("missing.json", "greeting")

		if err == nil || !strings.Contains(err.Error(), "document missing.json not found") {
			t.Fatalf("UpdateSegment() error = %v, want missing document context", err)
		}
		if len(engine.calls) != 0 {
			t.Errorf("engine calls = %d, want 0", len(engine.calls))
		}
	})

	t.Run("returns source segment errors without translating", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json")
		engine := &fakeEngine{result: translatedSegment}
		worker := NewWorker(testSettings(), project, engine, &fakeDocumentStore{}, &fakeReporter{})

		err := worker.UpdateSegment(source.Name, "missing")

		if !errors.Is(err, translation.ErrSegmentNotFound) {
			t.Fatalf("UpdateSegment() error = %v, want ErrSegmentNotFound", err)
		}
		if len(engine.calls) != 0 {
			t.Errorf("engine calls = %d, want 0", len(engine.calls))
		}
	})

	t.Run("stops at the destination whose translation fails", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json", translation.Segment{Key: "greeting", Value: "Hello"})
		translateErr := errors.New("translator unavailable")
		engine := &fakeEngine{result: func(segment translation.Segment, locale translation.Locale) (translation.Segment, error) {
			if locale == translation.Japanese {
				return translation.Segment{}, translateErr
			}
			return translatedSegment(segment, locale)
		}}
		worker := NewWorker(testSettings(), project, engine, &fakeDocumentStore{}, &fakeReporter{})

		err := worker.UpdateSegment(source.Name, "greeting")

		if !errors.Is(err, translateErr) || !strings.Contains(err.Error(), "from en to ja") {
			t.Fatalf("UpdateSegment() error = %v, want Japanese translation context", err)
		}
		assertSegment(t, project, "locales/fr/messages.json", translation.French, "greeting", "Hello in fr")
		if _, ok := project.Documents["locales/ja/messages.json"]; ok {
			t.Error("Japanese destination document exists after its translation failed")
		}
	})
}

func TestWorkerHandle(t *testing.T) {
	t.Run("processes inserts and persists every destination document", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json", translation.Segment{Key: "greeting", Value: "Hello"})
		store := &fakeDocumentStore{}
		worker := NewWorker(testSettings(), project, &fakeEngine{result: translatedSegment}, store, &fakeReporter{})

		err := worker.Handle(documentChanges(source.Name, sync.ChangeTypeInsert, "greeting"))

		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if got, want := documentNames(store.stored), []string{"locales/fr/messages.json", "locales/ja/messages.json"}; !reflect.DeepEqual(got, want) {
			t.Errorf("stored documents = %v, want %v", got, want)
		}
	})

	t.Run("creates all source segments and persists them", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json",
			translation.Segment{Key: "greeting", Value: "Hello"},
			translation.Segment{Key: "farewell", Value: "Goodbye"},
		)
		store := &fakeDocumentStore{}
		engine := &fakeEngine{result: translatedSegment}
		worker := NewWorker(testSettings(), project, engine, store, &fakeReporter{})

		err := worker.Handle(documentChanges(source.Name, sync.ChangeTypeCreateFile, ""))

		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if got, want := len(engine.calls), 4; got != want {
			t.Errorf("engine calls = %d, want %d", got, want)
		}
		assertSegment(t, project, "locales/fr/messages.json", translation.French, "farewell", "Goodbye in fr")
		if got := len(store.stored); got != 2 {
			t.Errorf("stored documents = %d, want 2", got)
		}
	})

	t.Run("deletes segments then persists target documents", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json", translation.Segment{Key: "greeting", Value: "Hello"})
		for _, locale := range testSettings().DestinationLocales {
			document := translation.NewDocument(source.DocumentNameForLocale(locale))
			mustSetSegment(t, document, locale, translation.Segment{Key: "greeting", Value: "translated"})
			mustSetSegment(t, document, locale, translation.Segment{Key: "kept", Value: "keep"})
			project.AddDocument(document)
		}
		store := &fakeDocumentStore{}

		err := NewWorker(testSettings(), project, &fakeEngine{}, store, &fakeReporter{}).Handle(documentChanges(source.Name, sync.ChangeTypeDelete, "greeting"))

		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		for _, locale := range testSettings().DestinationLocales {
			document, _ := project.GetDocumentIn(locale, source.Name)
			if _, err := document.GetSegment(locale, "greeting"); !errors.Is(err, translation.ErrSegmentNotFound) {
				t.Errorf("%s greeting lookup error = %v, want ErrSegmentNotFound", locale, err)
			}
		}
		if got := len(store.stored); got != 2 {
			t.Errorf("stored documents = %d, want 2", got)
		}
	})

	t.Run("deletes a document", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json")
		store := &fakeDocumentStore{}

		err := NewWorker(testSettings(), project, &fakeEngine{}, store, &fakeReporter{}).Handle(documentChanges(source.Name, sync.ChangeTypeDeleteFile, ""))

		if err != nil {
			t.Fatalf("Handle() error = %v", err)
		}
		if got, want := store.deleted, []string{"locales/fr/messages.json", "locales/ja/messages.json"}; !reflect.DeepEqual(got, want) {
			t.Errorf("deleted documents = %v, want %v", got, want)
		}
		if len(store.stored) != 0 {
			t.Errorf("stored documents = %d, want 0", len(store.stored))
		}
	})

	t.Run("rejects unknown changes without side effects", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json")
		store := &fakeDocumentStore{}

		err := NewWorker(testSettings(), project, &fakeEngine{}, store, &fakeReporter{}).Handle(documentChanges(source.Name, sync.ChangeType(99), ""))

		if err == nil || !strings.Contains(err.Error(), "unknown change type 99") {
			t.Fatalf("Handle() error = %v, want unknown change type", err)
		}
		if len(store.stored) != 0 || len(store.deleted) != 0 {
			t.Errorf("store side effects = stored %d, deleted %d, want none", len(store.stored), len(store.deleted))
		}
	})
}

func TestWorkerErrors(t *testing.T) {
	t.Run("remove segment stops when a target document is missing", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json")
		french := translation.NewDocument(source.DocumentNameForLocale(translation.French))
		mustSetSegment(t, french, translation.French, translation.Segment{Key: "greeting", Value: "Bonjour"})
		project.AddDocument(french)

		err := NewWorker(testSettings(), project, &fakeEngine{}, &fakeDocumentStore{}, &fakeReporter{}).RemoveSegment(source.Name, "greeting")

		if !errors.Is(err, translation.ErrDocumentNotFound) {
			t.Fatalf("RemoveSegment() error = %v, want ErrDocumentNotFound", err)
		}
		if _, err := french.GetSegment(translation.French, "greeting"); !errors.Is(err, translation.ErrSegmentNotFound) {
			t.Errorf("French segment lookup error = %v, want ErrSegmentNotFound", err)
		}
	})

	t.Run("wraps document store deletion errors", func(t *testing.T) {
		storeErr := errors.New("disk failure")
		store := &fakeDocumentStore{deleteErr: storeErr}

		err := NewWorker(testSettings(), translation.NewProject(), &fakeEngine{}, store, &fakeReporter{}).RemoveDocument("locales/en/messages.json")

		if !errors.Is(err, storeErr) || !strings.Contains(err.Error(), "locale fr") {
			t.Fatalf("RemoveDocument() error = %v, want French deletion context", err)
		}
		if got := len(store.deleted); got != 1 {
			t.Errorf("delete calls = %d, want 1", got)
		}
	})

	t.Run("wraps persistence errors and stops immediately", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json")
		for _, locale := range testSettings().DestinationLocales {
			document := translation.NewDocument(source.DocumentNameForLocale(locale))
			mustSetSegment(t, document, locale, translation.Segment{Key: "greeting", Value: "translated"})
			project.AddDocument(document)
		}
		storeErr := errors.New("disk failure")
		store := &fakeDocumentStore{storeErr: storeErr}

		err := NewWorker(testSettings(), project, &fakeEngine{}, store, &fakeReporter{}).persistChanges(source.Name)

		if !errors.Is(err, storeErr) || !strings.Contains(err.Error(), "locale fr") {
			t.Fatalf("persistChanges() error = %v, want French persistence context", err)
		}
		if got := len(store.stored); got != 1 {
			t.Errorf("store calls = %d, want 1", got)
		}
	})

	t.Run("returns missing destination context before storing", func(t *testing.T) {
		project, source := projectWithSource(t, "locales/en/messages.json")
		store := &fakeDocumentStore{}

		err := NewWorker(testSettings(), project, &fakeEngine{}, store, &fakeReporter{}).persistChanges(source.Name)

		if !errors.Is(err, translation.ErrDocumentNotFound) {
			t.Fatalf("persistChanges() error = %v, want ErrDocumentNotFound", err)
		}
		if len(store.stored) != 0 {
			t.Errorf("store calls = %d, want 0", len(store.stored))
		}
	})
}

func testSettings() *settings.Settings {
	return &settings.Settings{SourceLocale: translation.English, DestinationLocales: []translation.Locale{translation.French, translation.Japanese}}
}

func projectWithSource(t *testing.T, name string, segments ...translation.Segment) (*translation.Project, *translation.Document) {
	t.Helper()
	project := translation.NewProject()
	document := translation.NewDocument(name)
	for _, segment := range segments {
		mustSetSegment(t, document, translation.English, segment)
	}
	if document.Catalog(translation.English) == nil {
		if err := document.AddCatalog(translation.NewCatalog(translation.English)); err != nil {
			t.Fatalf("AddCatalog() error = %v", err)
		}
	}
	project.AddDocument(document)
	return project, document
}

func translatedSegment(segment translation.Segment, locale translation.Locale) (translation.Segment, error) {
	segment.Value += " in " + locale.String()
	return segment, nil
}

func mustSetSegment(t *testing.T, document *translation.Document, locale translation.Locale, segment translation.Segment) {
	t.Helper()
	if err := document.SetSegment(locale, segment); err != nil {
		t.Fatalf("SetSegment() error = %v", err)
	}
}

func assertSegment(t *testing.T, project *translation.Project, documentName string, locale translation.Locale, key, wantValue string) {
	t.Helper()
	document, ok := project.Documents[documentName]
	if !ok {
		t.Fatalf("project document %q not found", documentName)
	}
	segment, err := document.GetSegment(locale, key)
	if err != nil {
		t.Fatalf("GetSegment(%s, %s) error = %v", locale, key, err)
	}
	if segment.Value != wantValue {
		t.Errorf("segment value = %q, want %q", segment.Value, wantValue)
	}
}

func documentChanges(name string, changeType sync.ChangeType, segment string) sync.DocumentChanges {
	document := &sync.Document{Name: name}
	change := &sync.Change{Type: changeType, Document: document}
	if segment != "" {
		change.Segment = &segment
	}
	return sync.DocumentChanges{Document: document, Changes: []*sync.Change{change}}
}

func documentNames(documents []*translation.Document) []string {
	names := make([]string, len(documents))
	for index, document := range documents {
		names[index] = document.Name
	}
	return names
}
