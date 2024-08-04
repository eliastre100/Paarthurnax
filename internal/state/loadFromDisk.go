package state

import (
	"Paarthurnax/internal/translation"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
)

func LoadFromDisk(path string, verbose bool) (*State, error) {
	state := &State{}

	e := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.Name() == "fr.yml" {
			if verbose {
				log.Println(fmt.Sprintf("Processing %s", path))
			}

			t, err := translation.Load(path)
			if err != nil {
				return errors.New(fmt.Sprintf("Failed to load translation: %s", err))
			}

			state.Files = append(state.Files, TranslationFile{
				Path:           path,
				SegmentsHashes: HashTranslation(t.FlattenedSegments()),
			})
		}
		return err
	})
	if e != nil {
		return nil, errors.New("Unable to process the filestructure: " + e.Error())
	}

	return state, nil
}

func HashTranslation(m map[string]string) map[string]string {
	result := make(map[string]string)

	for k, v := range m {
		h := sha1.New()
		h.Write([]byte(v))
		result[k] = hex.EncodeToString(h.Sum(nil))
	}

	return result
}
