package document

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	projectadapter "Paarthurnax/internal/adapters/outbound/project"
	"Paarthurnax/internal/domain/translation"
)

func TestStorePersistsDocument(t *testing.T) {
	root := t.TempDir()
	name := filepath.Join("locales", "messages.yml")
	path := filepath.Join(root, name)
	document := translation.NewDocument(name)
	english := translation.NewCatalog(translation.English)
	english.AddSegment(*translation.NewSegment("greeting", "Hello"))
	english.AddSegment(*translation.NewSegment("menu.file", "File"))
	french := translation.NewCatalog(translation.French)
	french.AddSegment(*translation.NewSegment("greeting", "Bonjour"))
	french.AddSegment(*translation.NewSegment("menu.file", "Fichier"))
	for _, catalog := range []*translation.Catalog{english, french} {
		if err := document.AddCatalog(catalog); err != nil {
			t.Fatalf("AddCatalog() error = %v", err)
		}
	}

	if err := NewStore(root).Store(document); err != nil {
		t.Fatalf("Store() error = %v, want nil", err)
	}

	loaded, err := projectadapter.NewLoader(root).Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	got := loaded.Documents[path]
	if got == nil {
		t.Fatalf("stored document missing at %q", path)
	}
	assertStoredSegment(t, got, translation.English, "greeting", "Hello")
	assertStoredSegment(t, got, translation.English, "menu.file", "File")
	assertStoredSegment(t, got, translation.French, "greeting", "Bonjour")
	assertStoredSegment(t, got, translation.French, "menu.file", "Fichier")
}

func TestStoreRejectsConflictingSegmentKeys(t *testing.T) {
	document := translation.NewDocument(filepath.Join(t.TempDir(), "messages.yml"))
	catalog := translation.NewCatalog(translation.English)
	catalog.AddSegment(*translation.NewSegment("menu", "Menu"))
	catalog.AddSegment(*translation.NewSegment("menu.file", "File"))
	if err := document.AddCatalog(catalog); err != nil {
		t.Fatalf("AddCatalog() error = %v", err)
	}

	if err := NewStore("").Store(document); err == nil {
		t.Fatal("Store() error = nil, want conflicting key error")
	}
}

func TestStoreRejectsNilDocument(t *testing.T) {
	if err := NewStore("").Store(nil); err == nil {
		t.Fatal("Store() error = nil, want nil document error")
	}
}

func TestDeleteRemovesDocument(t *testing.T) {
	root := t.TempDir()
	name := "messages.yml"
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, []byte("en: {}\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := NewStore(root).Delete(name); err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Stat() error = %v, want %v", err, os.ErrNotExist)
	}
}

func TestDeleteIgnoresMissingDocument(t *testing.T) {
	if err := NewStore(t.TempDir()).Delete("missing.yml"); err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}
}

func assertStoredSegment(t *testing.T, document *translation.Document, locale translation.Locale, key, want string) {
	t.Helper()
	segment, err := document.GetSegment(locale, key)
	if err != nil {
		t.Fatalf("GetSegment(%q, %q) error = %v", locale, key, err)
	}
	if segment.Value != want {
		t.Errorf("GetSegment(%q, %q).Value = %q, want %q", locale, key, segment.Value, want)
	}
}
