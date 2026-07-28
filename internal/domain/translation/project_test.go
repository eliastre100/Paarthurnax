package translation

import "testing"

func TestProjectAddDocument(t *testing.T) {
	t.Run("indexes a document for each of its catalog locales", func(t *testing.T) {
		project := NewProject()
		document := documentWithCatalogs(t, "messages.json", English, French)

		project.AddDocument(document)

		if len(project.Documents) != 2 {
			t.Fatalf("len(Documents) = %d, want 2", len(project.Documents))
		}
		for _, locale := range []Locale{English, French} {
			if got := project.Documents[locale][document.Name]; got != document {
				t.Errorf("Documents[%q][%q] = %p, want %p", locale, document.Name, got, document)
			}
		}
	})

	t.Run("keeps documents with different names in the same locale", func(t *testing.T) {
		project := NewProject()
		messages := documentWithCatalogs(t, "messages.json", English)
		errors := documentWithCatalogs(t, "errors.json", English)

		project.AddDocument(messages)
		project.AddDocument(errors)

		documents := project.Documents[English]
		if len(documents) != 2 {
			t.Fatalf("len(Documents[%q]) = %d, want 2", English, len(documents))
		}
		if got := documents[messages.Name]; got != messages {
			t.Errorf("Documents[%q][%q] = %p, want %p", English, messages.Name, got, messages)
		}
		if got := documents[errors.Name]; got != errors {
			t.Errorf("Documents[%q][%q] = %p, want %p", English, errors.Name, got, errors)
		}
	})

	t.Run("does not create an index for a document without catalogs", func(t *testing.T) {
		project := NewProject()
		document := NewDocument("messages.json")

		project.AddDocument(&document)

		if len(project.Documents) != 0 {
			t.Errorf("len(Documents) = %d, want 0", len(project.Documents))
		}
	})
}

func documentWithCatalogs(t *testing.T, name string, locales ...Locale) *Document {
	t.Helper()

	document := NewDocument(name)
	for _, locale := range locales {
		if err := document.AddCatalog(&Catalog{Locale: locale}); err != nil {
			t.Fatalf("AddCatalog(%q) error = %v, want nil", locale, err)
		}
	}
	return &document
}
