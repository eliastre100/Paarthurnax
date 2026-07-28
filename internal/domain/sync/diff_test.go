package sync

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	t.Run("detects a single add action", func(t *testing.T) {
		previous := newSnapshotWithDocuments(t,
			newDocumentWithSegments(t, "messages.json", segmentContent{"greeting", "Hello"}),
		)
		currentDocument := newDocumentWithSegments(t, "messages.json",
			segmentContent{"greeting", "Hello"},
			segmentContent{"farewell", "Goodbye"},
		)
		current := newSnapshotWithDocuments(t, currentDocument)

		changes := Diff(previous, current)

		assertChanges(t, changes, []DocumentChanges{expectedDocumentChanges(currentDocument,
			Change{Type: ChangeTypeInsert, Document: &currentDocument, Segment: stringPointer("farewell")},
		)})
	})

	t.Run("detects a single remove action", func(t *testing.T) {
		previous := newSnapshotWithDocuments(t, newDocumentWithSegments(t, "messages.json",
			segmentContent{"greeting", "Hello"},
			segmentContent{"farewell", "Goodbye"},
		))
		currentDocument := newDocumentWithSegments(t, "messages.json", segmentContent{"greeting", "Hello"})
		current := newSnapshotWithDocuments(t, currentDocument)

		changes := Diff(previous, current)

		assertChanges(t, changes, []DocumentChanges{expectedDocumentChanges(currentDocument,
			Change{Type: ChangeTypeDelete, Document: &currentDocument, Segment: stringPointer("farewell")},
		)})
	})

	t.Run("detects a single update action", func(t *testing.T) {
		previous := newSnapshotWithDocuments(t,
			newDocumentWithSegments(t, "messages.json", segmentContent{"greeting", "Hello"}),
		)
		currentDocument := newDocumentWithSegments(t, "messages.json", segmentContent{"greeting", "Hi"})
		current := newSnapshotWithDocuments(t, currentDocument)

		changes := Diff(previous, current)

		assertChanges(t, changes, []DocumentChanges{expectedDocumentChanges(currentDocument,
			Change{Type: ChangeTypeUpdate, Document: &currentDocument, Segment: stringPointer("greeting")},
		)})
	})

	t.Run("detects a mixture of add, remove and update actions", func(t *testing.T) {
		previousDocument := newDocumentWithSegments(t, "messages.json",
			segmentContent{"unchanged", "Same"},
			segmentContent{"removed", "Remove me"},
			segmentContent{"updated", "Before"},
		)
		previous := newSnapshotWithDocuments(t, previousDocument)
		currentDocument := newDocumentWithSegments(t, "messages.json",
			segmentContent{"unchanged", "Same"},
			segmentContent{"updated", "After"},
			segmentContent{"added", "Add me"},
		)
		current := newSnapshotWithDocuments(t, currentDocument)

		changes := Diff(previous, current)

		assertChanges(t, changes, []DocumentChanges{expectedDocumentChanges(currentDocument,
			Change{Type: ChangeTypeDelete, Document: &currentDocument, Segment: stringPointer("removed")},
			Change{Type: ChangeTypeUpdate, Document: &currentDocument, Segment: stringPointer("updated")},
			Change{Type: ChangeTypeInsert, Document: &currentDocument, Segment: stringPointer("added")},
		)})
	})

	t.Run("detects a new document addition", func(t *testing.T) {
		previous := NewSnapshot("")
		currentDocument := newDocumentWithSegments(t, "new-messages.json", segmentContent{"greeting", "Hello"})
		current := newSnapshotWithDocuments(t, currentDocument)

		changes := Diff(previous, current)

		assertChanges(t, changes, []DocumentChanges{expectedDocumentChanges(currentDocument,
			Change{Type: ChangeTypeCreateFile, Document: &currentDocument},
		)})
	})

	t.Run("detects a document removal", func(t *testing.T) {
		previousDocument := newDocumentWithSegments(t, "old-messages.json", segmentContent{"greeting", "Hello"})
		previous := newSnapshotWithDocuments(t, previousDocument)
		current := NewSnapshot("")

		changes := Diff(previous, current)

		assertChanges(t, changes, []DocumentChanges{expectedDocumentChanges(previousDocument,
			Change{Type: ChangeTypeDeleteFile, Document: &previousDocument},
		)})
	})

	t.Run("detects an addition and a removal on move from a document to another", func(t *testing.T) {
		previousSource := newDocumentWithSegments(t, "source.json", segmentContent{"greeting", "Hello"})
		previousDestination := newDocumentWithSegments(t, "destination.json", segmentContent{"farewell", "Goodbye"})
		previous := newSnapshotWithDocuments(t, previousSource, previousDestination)
		currentSource := newDocumentWithSegments(t, "source.json")
		currentDestination := newDocumentWithSegments(t, "destination.json",
			segmentContent{"farewell", "Goodbye"},
			segmentContent{"greeting", "Hello"},
		)
		current := newSnapshotWithDocuments(t, currentSource, currentDestination)

		changes := Diff(previous, current)

		assertChanges(t, changes, []DocumentChanges{
			expectedDocumentChanges(currentDestination,
				Change{Type: ChangeTypeInsert, Document: &currentDestination, Segment: stringPointer("greeting")},
			),
			expectedDocumentChanges(currentSource,
				Change{Type: ChangeTypeDelete, Document: &currentSource, Segment: stringPointer("greeting")},
			),
		})
	})
}

type segmentContent struct {
	name    string
	content string
}

func newDocumentWithSegments(t *testing.T, name string, segments ...segmentContent) Document {
	t.Helper()

	document := NewDocument(name)
	for _, segment := range segments {
		digest := Digest(fmt.Sprintf("%x", sha1.Sum([]byte(segment.content))))
		if err := document.AddSegmentDigest(segment.name, digest); err != nil {
			t.Fatalf("AddSegmentDigest(%q) error = %v, want nil", segment.name, err)
		}
	}

	return *document
}

func newSnapshotWithDocuments(t *testing.T, documents ...Document) *Snapshot {
	t.Helper()

	snapshot := NewSnapshot("")
	for _, document := range documents {
		if err := snapshot.AddDocument(document); err != nil {
			t.Fatalf("AddDocument(%q) error = %v, want nil", document.Name, err)
		}
	}

	return snapshot
}

func stringPointer(value string) *string {
	return &value
}

func expectedDocumentChanges(document Document, changes ...Change) DocumentChanges {
	changePointers := make([]*Change, len(changes))
	for i := range changes {
		changePointers[i] = &changes[i]
	}
	return DocumentChanges{Document: &document, Changes: changePointers}
}

func assertChanges(t *testing.T, got, want []DocumentChanges) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Diff() = %#v, want %#v", got, want)
	}
}
