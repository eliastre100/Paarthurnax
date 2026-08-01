package project

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"Paarthurnax/internal/domain/translation"
)

func TestLoaderLoad(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "config.yml"), `
en:
  greeting: Hello
  menu:
    file: File
fr:
  greeting: Bonjour
  menu:
    file: Fichier
`)
	mustWriteFile(t, filepath.Join(root, "nested", "errors.yaml"), `
en:
  required: Required
`)
	mustWriteFile(t, filepath.Join(root, "ignored.txt"), "not yaml")
	mustWriteFile(t, filepath.Join(root, "uppercase.YAML"), "en:\n  uppercase: Loaded")

	project, err := NewLoader(root).Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if len(project.Documents) != 3 {
		t.Fatalf("len(Documents) = %d, want 3 documents", len(project.Documents))
	}

	configPath := filepath.Join(root, "config.yml")
	englishConfig := project.Documents[configPath]
	if englishConfig == nil {
		t.Fatalf("config document missing at %q", configPath)
	}
	if englishConfig != project.DocumentsWithCatalog(translation.English)[configPath] {
		t.Error("English catalog view should contain the indexed document")
	}
	if englishConfig != project.DocumentsWithCatalog(translation.French)[configPath] {
		t.Error("French catalog view should contain the indexed document")
	}
	if !englishConfig.Multilingual {
		t.Error("Multilingual = false, want true")
	}

	assertSegment(t, englishConfig.Catalogs[translation.English], "greeting", "Hello")
	assertSegment(t, englishConfig.Catalogs[translation.English], "menu.file", "File")
	assertSegment(t, englishConfig.Catalogs[translation.French], "greeting", "Bonjour")
	assertSegment(t, englishConfig.Catalogs[translation.French], "menu.file", "Fichier")

	errorsPath := filepath.Join(root, "nested", "errors.yaml")
	if document := project.DocumentsWithCatalog(translation.English)[errorsPath]; document == nil {
		t.Errorf("English errors document missing at %q", errorsPath)
	}

	uppercasePath := filepath.Join(root, "uppercase.YAML")
	uppercaseDocument := project.DocumentsWithCatalog(translation.English)[uppercasePath]
	if uppercaseDocument == nil {
		t.Fatalf("English uppercase document missing at %q", uppercasePath)
	}
	assertSegment(t, uppercaseDocument.Catalogs[translation.English], "uppercase", "Loaded")

	if got := len(project.DocumentsWithCatalog(translation.English)); got != 3 {
		t.Errorf("len(DocumentsWithCatalog(en)) = %d, want 3", got)
	}
}

func TestLoaderLoadErrors(t *testing.T) {
	t.Run("missing root", func(t *testing.T) {
		_, err := NewLoader(filepath.Join(t.TempDir(), "missing")).Load()

		if err == nil || !strings.Contains(err.Error(), "failed to load project") {
			t.Fatalf("Load() error = %v, want project loading error", err)
		}
	})

	t.Run("invalid YAML", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "invalid.yml"), "en: [")

		_, err := NewLoader(root).Load()

		if err == nil || !strings.Contains(err.Error(), "failed to unmarshal yaml") {
			t.Fatalf("Load() error = %v, want YAML unmarshalling error", err)
		}
	})

	t.Run("unsupported segment type", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "invalid.yml"), "en:\n  count: 1\n")

		_, err := NewLoader(root).Load()

		if err == nil || !strings.Contains(err.Error(), "unsupported type at key count: int") {
			t.Fatalf("Load() error = %v, want unsupported type error", err)
		}
	})
}

func TestLoadDocumentReadError(t *testing.T) {
	_, err := loadDocument(filepath.Join(t.TempDir(), "missing.yml"))

	if err == nil || !strings.Contains(err.Error(), "failed to read file") {
		t.Fatalf("loadDocument() error = %v, want file read error", err)
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("loadDocument() error = %v, want wrapped os.ErrNotExist", err)
	}
}

func mustWriteFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func assertSegment(t *testing.T, catalog *translation.Catalog, key, want string) {
	t.Helper()

	if catalog == nil {
		t.Fatal("catalog = nil")
	}
	segment, err := catalog.GetSegment(key)
	if err != nil {
		t.Fatalf("GetSegment(%q) error = %v", key, err)
	}
	if segment.Value != want {
		t.Errorf("GetSegment(%q).Value = %q, want %q", key, segment.Value, want)
	}
}
