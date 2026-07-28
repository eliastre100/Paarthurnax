package tomlproject

import (
	"fmt"

	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"

	"github.com/pelletier/go-toml/v2"
)

type stateV0 struct {
	Files []stateV0File `toml:"Files"`
}

type stateV0File struct {
	Path           string            `toml:"Path"`
	SegmentsHashes map[string]string `toml:"SegmentsHashes"`
}

func loadV0(data []byte) (*settings.ProjectState, error) {
	var persisted stateV0
	if err := toml.Unmarshal(data, &persisted); err != nil {
		return nil, fmt.Errorf("failed to decode project state version 0: %w", err)
	}

	snapshot := sync.NewSnapshot("fr")
	for _, file := range persisted.Files {
		if err := addDocument(snapshot, file.Path, file.SegmentsHashes); err != nil {
			return nil, err
		}
	}
	return &settings.ProjectState{
		Settings: &settings.Settings{
			SourceLocale: translation.French,
			DestinationLocales: []translation.Locale{
				"es", translation.English, "de", "it", "hu", "uk", "pl", "pt", "ro",
			},
		},
		Snapshot: snapshot,
	}, nil
}
