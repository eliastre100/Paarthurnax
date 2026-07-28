package sync

import (
	"Paarthurnax/internal/domain/translation"
	"crypto/sha1"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNewSnapshot(t *testing.T) {
	snapshot := NewSnapshot("en")

	if snapshot == nil {
		t.Fatal("NewSnapshot() = nil, want snapshot")
	}
	if snapshot.Locale != "en" {
		t.Errorf("Locale = %q, want en", snapshot.Locale)
	}
	if snapshot.Documents == nil {
		t.Fatal("Documents = nil, want initialized map")
	}
	if len(snapshot.Documents) != 0 {
		t.Errorf("len(Documents) = %d, want 0", len(snapshot.Documents))
	}
}

func TestNewSnapshotFromProject(t *testing.T) {
	project := translation.NewProject()
	messagesDocument := translation.NewDocument("messages.json")

	english := translation.NewCatalog(translation.English)
	english.AddSegment(*translation.NewSegment("greeting", "Hello"))
	english.AddSegment(*translation.NewSegment("farewell", "Goodbye"))

	french := translation.NewCatalog(translation.French)
	french.AddSegment(*translation.NewSegment("greeting", "Bonjour"))

	if err := messagesDocument.AddCatalog(english); err != nil {
		t.Fatalf("AddCatalog(English) error = %v, want nil", err)
	}
	if err := messagesDocument.AddCatalog(french); err != nil {
		t.Fatalf("AddCatalog(French) error = %v, want nil", err)
	}

	project.AddDocument(&messagesDocument)
	errorsDocument := translation.NewDocument("errors.json")
	frenchErrors := translation.NewCatalog(translation.French)
	frenchErrors.AddSegment(*translation.NewSegment("not_found", "Introuvable"))
	if err := errorsDocument.AddCatalog(frenchErrors); err != nil {
		t.Fatalf("AddCatalog(French) error = %v, want nil", err)
	}
	project.AddDocument(&errorsDocument)

	snapshot, err := NewSnapshotFromProject(project, translation.English)

	if err != nil {
		t.Fatalf("NewSnapshotFromProject() error = %v, want nil", err)
	}
	if snapshot.Locale != string(translation.English) {
		t.Errorf("Locale = %q, want %q", snapshot.Locale, translation.English)
	}
	if len(snapshot.Documents) != 1 {
		t.Fatalf("len(Documents) = %d, want 1", len(snapshot.Documents))
	}
	document, ok := snapshot.Documents["messages.json"]
	if !ok {
		t.Fatal("Documents[\"messages.json\"] is missing")
	}
	if got, want := document.GetSegmentDigest("greeting"), digestFor("Hello"); got != want {
		t.Errorf("GetSegmentDigest(\"greeting\") = %q, want %q", got, want)
	}
	if got, want := document.GetSegmentDigest("farewell"), digestFor("Goodbye"); got != want {
		t.Errorf("GetSegmentDigest(\"farewell\") = %q, want %q", got, want)
	}
	if _, ok := snapshot.Documents["errors.json"]; ok {
		t.Error("Documents contains a document that has no English catalog")
	}
}

func TestNewSnapshotFromProjectMissingLocale(t *testing.T) {
	project := translation.NewProject()

	snapshot, err := NewSnapshotFromProject(project, translation.Japanese)

	if err != nil {
		t.Fatalf("NewSnapshotFromProject() error = %v, want nil", err)
	}
	if snapshot.Locale != string(translation.Japanese) {
		t.Errorf("Locale = %q, want %q", snapshot.Locale, translation.Japanese)
	}
	if len(snapshot.Documents) != 0 {
		t.Errorf("len(Documents) = %d, want 0", len(snapshot.Documents))
	}
}

func TestSnapshotAddDocument(t *testing.T) {
	t.Run("adds documents by name", func(t *testing.T) {
		snapshot := NewSnapshot("en")
		documents := []Document{
			newDocumentWithSegment(t, "en.messages.foo", "greeting", "Hello"),
			newDocumentWithSegment(t, "en.messages.bar", "greeting", "Bonjour"),
		}

		for _, document := range documents {
			if err := snapshot.AddDocument(document); err != nil {
				t.Fatalf("AddDocument(%q) error = %v, want nil", document.Name, err)
			}
		}

		if len(snapshot.Documents) != len(documents) {
			t.Fatalf("len(Documents) = %d, want %d", len(snapshot.Documents), len(documents))
		}
		for _, document := range documents {
			if got := snapshot.Documents[document.Name]; !reflect.DeepEqual(got, document) {
				t.Errorf("Documents[%q] = %#v, want %#v", document.Name, got, document)
			}
		}
	})

	t.Run("rejects an existing document with the same name", func(t *testing.T) {
		snapshot := NewSnapshot("en")
		original := newDocumentWithSegment(t, "messages.en", "greeting", "Hello")
		duplicate := newDocumentWithSegment(t, original.Name, "greeting", "Hi")

		if err := snapshot.AddDocument(original); err != nil {
			t.Fatalf("AddDocument() first error = %v, want nil", err)
		}

		err := snapshot.AddDocument(duplicate)

		if !errors.Is(err, ErrDocumentAlreadyExists) {
			t.Fatalf("AddDocument() error = %v, want %v", err, ErrDocumentAlreadyExists)
		}
		if len(snapshot.Documents) != 1 {
			t.Fatalf("len(Documents) = %d, want 1", len(snapshot.Documents))
		}
		if got := snapshot.Documents[original.Name]; !reflect.DeepEqual(got, original) {
			t.Errorf("Documents[%q] = %#v, want original document %#v", original.Name, got, original)
		}
	})
}

func newDocumentWithSegment(t *testing.T, name, segment, content string) Document {
	t.Helper()

	document := NewDocument(name)
	digest := Digest(fmt.Sprintf("%x", sha1.Sum([]byte(content))))
	if err := document.AddSegmentDigest(segment, digest); err != nil {
		t.Fatalf("AddSegmentDigest(%q) error = %v, want nil", segment, err)
	}

	return *document
}

func digestFor(value string) Digest {
	return Digest(fmt.Sprintf("%x", sha1.Sum([]byte(value))))
}
