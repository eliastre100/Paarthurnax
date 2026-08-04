package normalize

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	documentadapter "Paarthurnax/internal/adapters/outbound/document"
	projectadapter "Paarthurnax/internal/adapters/outbound/project"
	"Paarthurnax/internal/domain/translation"
)

func TestExecuteNormalizesYAMLKeyOrder(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "config.yml")
	const input = `fr:
  zoo: Zoo
  menu:
    save: Sauvegarder
    cancel: Annuler
  apple: Pomme
en:
  zoo: Zoo
  menu:
    save: Save
    cancel: Cancel
  apple: Apple
`
	if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	secondPath := filepath.Join(root, "errors.yaml")
	if err := os.WriteFile(secondPath, []byte("en:\n  zebra: Zebra\n  apple: Apple\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", secondPath, err)
	}

	if err := Execute(projectadapter.NewLoader(root), documentadapter.NewStore("")); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	got := string(data)
	for _, orderedKeys := range [][]string{
		{"en:", "fr:"},
		{"    apple:", "    menu:", "    zoo:"},
		{"        cancel:", "        save:"},
	} {
		assertOrdered(t, got, orderedKeys...)
	}

	secondData, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", secondPath, err)
	}
	assertOrdered(t, string(secondData), "  apple:", "  zebra:")
}

func TestExecuteReturnsLoadError(t *testing.T) {
	loadErr := errors.New("project unavailable")
	store := &recordingStore{}

	err := Execute(&stubLoader{err: loadErr}, store)

	if !errors.Is(err, loadErr) {
		t.Fatalf("Execute() error = %v, want wrapped %v", err, loadErr)
	}
	if !strings.Contains(err.Error(), "loading project") {
		t.Errorf("Execute() error = %q, want loading context", err)
	}
	if len(store.documents) != 0 {
		t.Errorf("stored documents = %d, want 0", len(store.documents))
	}
}

func TestExecuteReturnsStoreError(t *testing.T) {
	storeErr := errors.New("disk full")
	document := translation.NewDocument("config.yml")
	project := translation.NewProject()
	project.AddDocument(document)
	store := &recordingStore{err: storeErr}

	err := Execute(&stubLoader{project: project}, store)

	if !errors.Is(err, storeErr) {
		t.Fatalf("Execute() error = %v, want wrapped %v", err, storeErr)
	}
	if !strings.Contains(err.Error(), "storing document") {
		t.Errorf("Execute() error = %q, want storing context", err)
	}
	if len(store.documents) != 1 || store.documents[0] != document {
		t.Errorf("stored documents = %#v, want [%#v]", store.documents, document)
	}
}

func assertOrdered(t *testing.T, text string, keys ...string) {
	t.Helper()

	position := -1
	for _, key := range keys {
		next := strings.Index(text, key)
		if next == -1 {
			t.Fatalf("normalized YAML does not contain %q:\n%s", key, text)
		}
		if next < position {
			t.Errorf("key %q appears before the preceding key in:\n%s", key, text)
		}
		position = next
	}
}
