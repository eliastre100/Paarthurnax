package tomlproject

import (
	"fmt"
	"time"

	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"

	"github.com/pelletier/go-toml/v2"
)

type stateV1 struct {
	Version    int            `toml:"version"`
	Locales    stateV1Locales `toml:"locales"`
	LastUpdate time.Time      `toml:"last_update"`
	Files      []stateV1File  `toml:"files"`
}

type stateV1Locales struct {
	Source       string   `toml:"source"`
	Destinations []string `toml:"destinations"`
}

type stateV1File struct {
	Path     string            `toml:"path"`
	Segments map[string]string `toml:"segments"`
}

func loadV1(data []byte) (*settings.ProjectState, error) {
	var persisted stateV1
	if err := toml.Unmarshal(data, &persisted); err != nil {
		return nil, fmt.Errorf("failed to decode project state version 1: %w", err)
	}

	source := translation.Locale(persisted.Locales.Source)
	if source == "" {
		source = translation.French
	}

	snapshot := sync.NewSnapshot(string(source))
	for _, file := range persisted.Files {
		if err := addDocument(snapshot, file.Path, file.Segments); err != nil {
			return nil, err
		}
	}

	var destinations []translation.Locale
	if len(persisted.Locales.Destinations) == 0 {
		destinations = []translation.Locale{
			translation.Spanish, translation.English, translation.German, translation.Italian, translation.Hungarian, translation.Ukrainian, translation.Polish, translation.Portuguese, translation.Romanian,
		}
	} else {
		destinations = make([]translation.Locale, len(persisted.Locales.Destinations))
		for i, locale := range persisted.Locales.Destinations {
			destinations[i] = translation.Locale(locale)
		}
	}
	return &settings.ProjectState{
		Settings: &settings.Settings{
			SourceLocale:       source,
			DestinationLocales: destinations,
		},
		Snapshot: snapshot,
	}, nil
}
