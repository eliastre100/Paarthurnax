package sync

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"testing"
)

func TestNewDocument(t *testing.T) {
	document := NewDocument("en.messages.json")

	if document == nil {
		t.Fatal("NewDocument() = nil, want document")
	}
	if document.Name != "en.messages.json" {
		t.Errorf("Name = %q, want %q", document.Name, "en.messages.json")
	}
	if document.Digests == nil {
		t.Fatal("Digests = nil, want initialized map")
	}
	if len(document.Digests) != 0 {
		t.Errorf("len(Digests) = %d, want 0", len(document.Digests))
	}
}

func TestDocumentAddAndGetSegmentDigest(t *testing.T) {
	document := NewDocument("en.messages.json")
	content := "Hello, world!"
	digest := Digest(fmt.Sprintf("%x", sha1.Sum([]byte(content))))

	if err := document.AddSegmentDigest("greeting", digest); err != nil {
		t.Fatalf("AddSegmentDigest() error = %v, want nil", err)
	}

	if got := document.GetSegmentDigest("greeting"); got != digest {
		t.Errorf("GetSegmentDigest(%q) = %q, want SHA-1 digest %q", "greeting", got, digest)
	}
}

func TestDocumentAddSegmentDigestRejectsExistingSegment(t *testing.T) {
	document := NewDocument("en.messages.json")
	segment := "greeting"
	original := Digest(fmt.Sprintf("%x", sha1.Sum([]byte("Hello"))))
	replacement := Digest(fmt.Sprintf("%x", sha1.Sum([]byte("Hi"))))

	if err := document.AddSegmentDigest(segment, original); err != nil {
		t.Fatalf("AddSegmentDigest() first error = %v, want nil", err)
	}

	err := document.AddSegmentDigest(segment, replacement)

	if !errors.Is(err, ErrSegmentAlreadyExists) {
		t.Fatalf("AddSegmentDigest() error = %v, want %v", err, ErrSegmentAlreadyExists)
	}
	if len(document.Digests) != 1 {
		t.Fatalf("len(Digests) = %d, want 1", len(document.Digests))
	}
	if got := document.GetSegmentDigest(segment); got != original {
		t.Errorf("GetSegmentDigest(%q) = %q, want original SHA-1 digest %q", segment, got, original)
	}
}

func TestDocumentGetSegmentDigestMissing(t *testing.T) {
	document := NewDocument("en.messages.json")

	if got := document.GetSegmentDigest("missing"); got != "" {
		t.Errorf("GetSegmentDigest(%q) = %q, want empty digest", "missing", got)
	}
}
