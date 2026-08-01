package translation

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewCatalog(t *testing.T) {
	catalog := NewCatalog("fr")

	if _, err := uuid.Parse(catalog.ID); err != nil {
		t.Errorf("ID = %q, want UUID: %v", catalog.ID, err)
	}
	if catalog.Locale != "fr" {
		t.Errorf("Locale = %q, want %q", catalog.Locale, "fr")
	}
	if catalog.Segments == nil {
		t.Fatal("Segments = nil, want initialized map")
	}
	if len(catalog.Segments) != 0 {
		t.Errorf("len(Segments) = %d, want 0", len(catalog.Segments))
	}
}

func TestCatalogAddAndGetSegment(t *testing.T) {
	catalog := NewCatalog("en")
	segment := Segment{Key: "greeting", Value: "Hello"}

	catalog.AddSegment(segment)

	got, err := catalog.GetSegment(segment.Key)
	if err != nil {
		t.Fatalf("GetSegment() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, segment) {
		t.Errorf("GetSegment() = %#v, want %#v", got, segment)
	}
}

func TestCatalogAddSegmentReplacesExistingSegment(t *testing.T) {
	catalog := NewCatalog("en")
	original := Segment{Key: "greeting", Value: "Hello"}
	updated := Segment{Key: "greeting", Value: "Hi", Plural: true}
	catalog.AddSegment(original)

	catalog.AddSegment(updated)

	if len(catalog.Segments) != 1 {
		t.Fatalf("len(Segments) = %d, want 1", len(catalog.Segments))
	}
	if got := catalog.Segments[updated.Key]; !reflect.DeepEqual(got, updated) {
		t.Errorf("Segments[%q] = %#v, want %#v", updated.Key, got, updated)
	}
}

func TestCatalogGetSegmentNotFound(t *testing.T) {
	catalog := NewCatalog("en")

	got, err := catalog.GetSegment("missing")

	if !errors.Is(err, ErrSegmentNotFound) {
		t.Fatalf("GetSegment() error = %v, want %v", err, ErrSegmentNotFound)
	}
	if !reflect.DeepEqual(got, Segment{}) {
		t.Errorf("GetSegment() = %#v, want zero Segment", got)
	}
}

func TestCatalogUpdateSegment(t *testing.T) {
	catalog := NewCatalog("en")
	original := Segment{Key: "greeting", Value: "Hello"}
	updated := Segment{Key: "greeting", Value: "Hi", Plural: true}
	catalog.AddSegment(original)

	err := catalog.UpdateSegment(original.Key, updated)

	if err != nil {
		t.Fatalf("UpdateSegment() error = %v, want nil", err)
	}
	if got := catalog.Segments[original.Key]; !reflect.DeepEqual(got, updated) {
		t.Errorf("Segments[%q] = %#v, want %#v", original.Key, got, updated)
	}
}

func TestCatalogUpdateSegmentNotFound(t *testing.T) {
	catalog := NewCatalog("en")
	segment := Segment{Key: "greeting", Value: "Hello"}

	err := catalog.UpdateSegment(segment.Key, segment)

	if !errors.Is(err, ErrSegmentNotFound) {
		t.Fatalf("UpdateSegment() error = %v, want %v", err, ErrSegmentNotFound)
	}
	if len(catalog.Segments) != 0 {
		t.Errorf("len(Segments) = %d, want 0", len(catalog.Segments))
	}
}

func TestCatalogRemoveSegment(t *testing.T) {
	catalog := NewCatalog("en")
	segment := Segment{Key: "greeting", Value: "Hello"}
	catalog.AddSegment(segment)

	catalog.RemoveSegment(segment.Key)

	if _, ok := catalog.Segments[segment.Key]; ok {
		t.Errorf("Segments[%q] still exists after RemoveSegment", segment.Key)
	}
}

func TestCatalogIsEmpty(t *testing.T) {
	catalog := NewCatalog("en")

	if !catalog.isEmpty() {
		t.Error("isEmpty() = false, want true for a new catalog")
	}

	segment := Segment{Key: "greeting", Value: "Hello"}
	catalog.AddSegment(segment)
	if catalog.isEmpty() {
		t.Error("isEmpty() = true, want false after adding a segment")
	}

	catalog.RemoveSegment(segment.Key)
	if !catalog.isEmpty() {
		t.Error("isEmpty() = false, want true after removing the final segment")
	}
}
