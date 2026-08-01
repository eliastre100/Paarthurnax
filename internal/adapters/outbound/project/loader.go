package project

import (
	"Paarthurnax/internal/domain/translation"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	Root string
}

func NewLoader(root string) *Loader {
	return &Loader{Root: root}
}

func (l *Loader) Load() (*translation.Project, error) {
	project := translation.NewProject()

	if err := filepath.WalkDir(l.Root, func(path string, d fs.DirEntry, err error) error {
		ext := strings.ToLower(filepath.Ext(path))
		if err == nil && (ext == ".yml" || ext == ".yaml") {
			document, err := loadDocument(path)
			if err != nil {
				return fmt.Errorf("failed to load document %s: %w", path, err)
			}

			project.AddDocument(document)
		}
		return err
	}); err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	return project, nil
}

func loadDocument(path string) (*translation.Document, error) {
	var yml map[string]map[string]any
	document := translation.NewDocument(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &yml); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	for locale, segments := range yml {
		catalog := translation.NewCatalog(translation.Locale(locale))
		err := fillCatalog(catalog, segments, "")
		if err != nil {
			return nil, fmt.Errorf("failed to read segments: %w", err)
		}
		if err := document.AddCatalog(catalog); err != nil {
			return nil, fmt.Errorf("failed to add catalog: %w", err)
		}
	}

	if len(document.Catalogs) > 1 {
		document.Multilingual = true
	}

	return document, nil
}

func fillCatalog(catalog *translation.Catalog, segments map[string]any, prefix string) error {
	for key, value := range segments {
		segmentKey := key
		if prefix != "" {
			segmentKey = prefix + "." + key
		}

		switch value.(type) {
		case map[string]any:
			if err := fillCatalog(catalog, value.(map[string]any), segmentKey); err != nil {
				return err
			}
		case string:
			catalog.AddSegment(*translation.NewSegment(segmentKey, value.(string)))
		default:
			return fmt.Errorf("unsupported type at key %s: %T", segmentKey, value)
		}
	}
	return nil
}
