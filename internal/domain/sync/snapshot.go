package sync

import (
	"Paarthurnax/internal/domain/translation"
	"crypto/sha1"
	"errors"
	"fmt"
)

type Snapshot struct {
	Locale    string
	Documents map[string]Document
}

var ErrDocumentAlreadyExists = errors.New("document already exists")

func NewSnapshot(locale string) *Snapshot {
	return &Snapshot{
		Locale:    locale,
		Documents: make(map[string]Document),
	}
}

func NewSnapshotFromProject(project *translation.Project, locale translation.Locale) (*Snapshot, error) {
	snapshot := NewSnapshot(string(locale))

	for _, sourceDocument := range project.Documents[locale] {
		document := NewDocument(sourceDocument.Name)
		catalog := sourceDocument.Catalogs[locale]

		for key, segment := range catalog.Segments {
			digest := Digest(fmt.Sprintf("%x", sha1.Sum([]byte(segment.Value))))
			if err := document.AddSegmentDigest(key, digest); err != nil {
				return nil, fmt.Errorf("failed to create snapshot document %q: %w", document.Name, err)
			}
		}

		if err := snapshot.AddDocument(*document); err != nil {
			return nil, fmt.Errorf("failed to add snapshot document %q: %w", document.Name, err)
		}
	}

	return snapshot, nil
}

func (s *Snapshot) AddDocument(doc Document) error {
	if _, ok := s.Documents[doc.Name]; ok {
		return ErrDocumentAlreadyExists
	}
	s.Documents[doc.Name] = doc
	return nil
}
