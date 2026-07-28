package tomlproject

import (
	"os"
	"path/filepath"
	"testing"

	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"
)

func newProjectState(t *testing.T) *settings.ProjectState {
	t.Helper()

	snapshot := sync.NewSnapshot("en")
	document := sync.NewDocument("config/locales/en.yml")
	if err := document.AddSegmentDigest("navigation.home", sync.Digest("abc123")); err != nil {
		t.Fatal(err)
	}
	if err := document.AddSegmentDigest("navigation.exit", sync.Digest("def456")); err != nil {
		t.Fatal(err)
	}
	if err := snapshot.AddDocument(*document); err != nil {
		t.Fatal(err)
	}

	return &settings.ProjectState{
		Settings: &settings.Settings{
			SourceLocale:       translation.English,
			DestinationLocales: []translation.Locale{translation.French, translation.Japanese},
		},
		Snapshot: snapshot,
	}
}

func loadState(t *testing.T, contents string) *settings.ProjectState {
	t.Helper()

	path := filepath.Join(t.TempDir(), ".paarthurnax")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	state, err := NewRepository(path).Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	return state
}

func assertDocumentDigest(t *testing.T, snapshot *sync.Snapshot, name, segment string, want sync.Digest) {
	t.Helper()

	document, ok := snapshot.Documents[name]
	if !ok {
		t.Fatalf("document %q missing", name)
	}
	if got := document.GetSegmentDigest(segment); got != want {
		t.Errorf("GetSegmentDigest(%q) = %q, want %q", segment, got, want)
	}
}
