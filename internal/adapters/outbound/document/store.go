package document

import (
	"Paarthurnax/internal/domain/translation"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Store struct {
	Root string
}

func NewStore(root string) *Store {
	return &Store{Root: root}
}

func (s *Store) Store(document *translation.Document) error {
	if document == nil {
		return fmt.Errorf("cannot store a nil document")
	}

	persisted := make(map[string]map[string]any, len(document.Catalogs))
	for locale, catalog := range document.Catalogs {
		if catalog == nil {
			return fmt.Errorf("cannot store document %s: catalog for locale %s is nil", document.Name, locale)
		}

		segments, err := nestSegments(catalog.Segments)
		if err != nil {
			return fmt.Errorf("cannot store document %s: failed to encode catalog %s: %w", document.Name, locale, err)
		}
		persisted[locale.String()] = segments
	}

	data, err := yaml.Marshal(persisted)
	if err != nil {
		return fmt.Errorf("failed to marshal document %s: %w", document.Name, err)
	}
	path := filepath.Join(s.Root, document.Name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create directory for document %s: %w", document.Name, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write document %s: %w", document.Name, err)
	}
	return nil
}

func (s *Store) Delete(document string) error {
	if err := os.Remove(filepath.Join(s.Root, document)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to delete document %s: %w", document, err)
	}
	return nil
}

func fillSegment(tree map[string]any, keyParts []string, value string) error {
	if len(keyParts) == 1 {
		tree[keyParts[0]] = value
		return nil
	}

	subtree, ok := tree[keyParts[0]]
	if !ok {
		subtree = make(map[string]any)
		tree[keyParts[0]] = subtree
	}

	castedSubtree, ok := subtree.(map[string]any)
	if !ok {
		return fmt.Errorf("segment key %q conflicts with nested segment", keyParts[0])
	}

	return fillSegment(castedSubtree, keyParts[1:], value)
}

func nestSegments(segments map[string]translation.Segment) (map[string]any, error) {
	root := make(map[string]any)

	for key, segment := range segments {
		parts := strings.Split(key, ".")
		if err := fillSegment(root, parts, segment.Value); err != nil {
			return nil, err
		}
	}
	return root, nil
}
