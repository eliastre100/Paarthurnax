package v0

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
)

type TranslationState struct {
	Files []TranslationFile
}

type TranslationFile struct {
	Path           string
	SegmentsHashes map[string]string // Key: sha1
}

func Load(data []byte, version int) (*TranslationState, error) {
	var state TranslationState

	if version != 0 {
		return nil, fmt.Errorf("unsupported state version: %d", version)
	}

	if err := toml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unable to parse state file: %w", err)
	}
	return &state, nil
}
