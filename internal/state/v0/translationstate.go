package v0

import (
	"Paarthurnax/internal/translation"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type TranslationState struct {
	Files []TranslationFile
}

func (state *TranslationState) GetTranslationFile(path string) *TranslationFile {
	for _, file := range state.Files {
		if file.Path == path {
			return &file
		}
	}
	return nil
}

func (state *TranslationState) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	encoded, err := toml.Marshal(state)
	if err != nil {
		return err
	}
	if _, err = f.Write(encoded); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

func LoadState(path string) (*TranslationState, error) {
	var state TranslationState
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open state file: %w", err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			slog.Error(fmt.Sprintf("Failed to close state paarthurnax state file %s: %w", path, err))
		}
	}(f)
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read state file: %w", err)
	}
	err = toml.Unmarshal(data, &state)
	if err != nil {
		return nil, fmt.Errorf("unable to parse state file: %w", err)
	}
	return &state, nil
}

func BuildFromDisk(path string, verbose bool) (*TranslationState, error) {
	state := &TranslationState{}

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.Name() == "fr.yml" {
			if verbose {
				slog.Info("Processing", "path", path)
			}

			t, err := translation.Load(path)
			if err != nil {
				return fmt.Errorf("failed to load translation: %w", err)
			}

			state.Files = append(state.Files, TranslationFile{
				Path:           path,
				SegmentsHashes: hashTranslation(t.FlattenedSegments()),
			})
		}
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("unable to process the filestructure: %w", err)
	}

	return state, nil
}

func hashTranslation(m map[string]string) map[string]string {
	result := make(map[string]string)

	for k, v := range m {
		h := sha1.New()
		h.Write([]byte(v))
		result[k] = hex.EncodeToString(h.Sum(nil))
	}

	return result
}
