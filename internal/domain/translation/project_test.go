package translation

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewProject(t *testing.T) {
	project := NewProject()

	if project.Documents == nil {
		t.Fatal("Documents = nil, want initialized map")
	}
	if len(project.Documents) != 0 {
		t.Errorf("len(Documents) = %d, want 0", len(project.Documents))
	}
}

func TestProjectAddDocument(t *testing.T) {
	t.Run("indexes a document by name", func(t *testing.T) {
		project := NewProject()
		document := documentWithCatalogs(t, "messages.json", English, French)

		project.AddDocument(document)

		if len(project.Documents) != 1 {
			t.Fatalf("len(Documents) = %d, want 1", len(project.Documents))
		}
		if got := project.Documents[document.Name]; got != document {
			t.Errorf("Documents[%q] = %p, want %p", document.Name, got, document)
		}
	})

	t.Run("keeps documents with different names", func(t *testing.T) {
		project := NewProject()
		messages := documentWithCatalogs(t, "messages.json", English)
		errorsDocument := documentWithCatalogs(t, "errors.json", English)

		project.AddDocument(messages)
		project.AddDocument(errorsDocument)

		if len(project.Documents) != 2 {
			t.Fatalf("len(Documents) = %d, want 2", len(project.Documents))
		}
		if got := project.Documents[messages.Name]; got != messages {
			t.Errorf("Documents[%q] = %p, want %p", messages.Name, got, messages)
		}
		if got := project.Documents[errorsDocument.Name]; got != errorsDocument {
			t.Errorf("Documents[%q] = %p, want %p", errorsDocument.Name, got, errorsDocument)
		}
	})

	t.Run("replaces a document with the same name", func(t *testing.T) {
		project := NewProject()
		first := documentWithCatalogs(t, "messages.json", English)
		replacement := documentWithCatalogs(t, "messages.json", French)

		project.AddDocument(first)
		project.AddDocument(replacement)

		if len(project.Documents) != 1 {
			t.Fatalf("len(Documents) = %d, want 1", len(project.Documents))
		}
		if got := project.Documents[replacement.Name]; got != replacement {
			t.Errorf("Documents[%q] = %p, want replacement %p", replacement.Name, got, replacement)
		}
	})
}

func TestProjectDocumentsWithCatalog(t *testing.T) {
	project := NewProject()
	englishOnly := documentWithCatalogs(t, "messages.json", English)
	frenchOnly := documentWithCatalogs(t, "errors.json", French)
	bilingual := documentWithCatalogs(t, "common.json", English, French)
	withoutCatalogs := NewDocument("empty.json")

	for _, document := range []*Document{englishOnly, frenchOnly, bilingual, withoutCatalogs} {
		project.AddDocument(document)
	}

	tests := []struct {
		name   string
		locale Locale
		want   map[string]*Document
	}{
		{
			name:   "returns documents with English catalogs",
			locale: English,
			want: map[string]*Document{
				englishOnly.Name: englishOnly,
				bilingual.Name:   bilingual,
			},
		},
		{
			name:   "returns documents with French catalogs",
			locale: French,
			want: map[string]*Document{
				frenchOnly.Name: frenchOnly,
				bilingual.Name:  bilingual,
			},
		},
		{name: "returns an empty initialized map when no documents match", locale: Japanese, want: map[string]*Document{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := project.DocumentsWithCatalog(tt.locale)

			if got == nil {
				t.Fatal("DocumentsWithCatalog() = nil, want initialized map")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DocumentsWithCatalog(%q) = %#v, want %#v", tt.locale, got, tt.want)
			}
		})
	}
}

func TestProjectGetDocumentIn(t *testing.T) {
	t.Run("returns the document when it already has the requested locale", func(t *testing.T) {
		project := NewProject()
		document := documentWithCatalogs(t, "messages.en.yml", English, French)
		project.AddDocument(document)

		got, err := project.GetDocumentIn(French, document.Name)

		if err != nil {
			t.Fatalf("GetDocumentIn() error = %v, want nil", err)
		}
		if got != document {
			t.Errorf("GetDocumentIn() = %p, want %p", got, document)
		}
	})

	t.Run("resolves the locale-specific document name", func(t *testing.T) {
		project := NewProject()
		source := documentWithCatalogs(t, "messages.en.yml", English)
		destination := documentWithCatalogs(t, "messages.fr.yml", French)
		project.AddDocument(source)
		project.AddDocument(destination)

		got, err := project.GetDocumentIn(French, source.Name)

		if err != nil {
			t.Fatalf("GetDocumentIn() error = %v, want nil", err)
		}
		if got != destination {
			t.Errorf("GetDocumentIn() = %p, want %p", got, destination)
		}
	})

	t.Run("returns not found when the source document does not exist", func(t *testing.T) {
		project := NewProject()

		got, err := project.GetDocumentIn(French, "missing.en.yml")

		if !errors.Is(err, ErrDocumentNotFound) {
			t.Fatalf("GetDocumentIn() error = %v, want %v", err, ErrDocumentNotFound)
		}
		if got != nil {
			t.Errorf("GetDocumentIn() = %p, want nil", got)
		}
	})

	t.Run("returns not found when the locale-specific document does not exist", func(t *testing.T) {
		project := NewProject()
		source := documentWithCatalogs(t, "messages.en.yml", English)
		project.AddDocument(source)

		got, err := project.GetDocumentIn(French, source.Name)

		if !errors.Is(err, ErrDocumentNotFound) {
			t.Fatalf("GetDocumentIn() error = %v, want %v", err, ErrDocumentNotFound)
		}
		if got != nil {
			t.Errorf("GetDocumentIn() = %p, want nil", got)
		}
	})
}

func documentWithCatalogs(t *testing.T, name string, locales ...Locale) *Document {
	t.Helper()

	document := NewDocument(name)
	for _, locale := range locales {
		if err := document.AddCatalog(NewCatalog(locale)); err != nil {
			t.Fatalf("AddCatalog(%q) error = %v, want nil", locale, err)
		}
	}
	return document
}
