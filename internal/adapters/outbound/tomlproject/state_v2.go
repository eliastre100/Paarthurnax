package tomlproject

import (
	"errors"
	"fmt"
	"time"

	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"

	"github.com/pelletier/go-toml/v2"
)

type stateV2 struct {
	Version    int            `toml:"version"`
	Locales    stateV2Locales `toml:"locales"`
	LastUpdate time.Time      `toml:"last_update"`
	Files      []stateV2File  `toml:"files"`
}

type stateV2Locales struct {
	Source       string   `toml:"source"`
	Destinations []string `toml:"destinations"`
}

type stateV2File struct {
	Path     string            `toml:"path"`
	Segments map[string]string `toml:"segments"`
}

func loadV2(data []byte) (*settings.ProjectState, error) {
	var persisted stateV2
	if err := toml.Unmarshal(data, &persisted); err != nil {
		return nil, fmt.Errorf("failed to decode project state version 2: %w", err)
	}

	snapshot := sync.NewSnapshot(persisted.Locales.Source)
	for _, file := range persisted.Files {
		if err := addDocument(snapshot, file.Path, file.Segments); err != nil {
			return nil, err
		}
	}

	destinations := make([]translation.Locale, len(persisted.Locales.Destinations))
	for i, locale := range persisted.Locales.Destinations {
		destinations[i] = translation.Locale(locale)
	}
	return &settings.ProjectState{
		Settings: &settings.Settings{
			SourceLocale:       translation.Locale(persisted.Locales.Source),
			DestinationLocales: destinations,
		},
		Snapshot: snapshot,
	}, nil
}

func encodeV2(state *settings.ProjectState) (*stateV2, error) {
	if state == nil {
		return nil, errors.New("project state cannot be nil")
	}
	if state.Settings == nil {
		return nil, errors.New("project settings cannot be nil")
	}
	if state.Snapshot == nil {
		return nil, errors.New("project snapshot cannot be nil")
	}

	destinations := make([]string, len(state.Settings.DestinationLocales))
	for i, locale := range state.Settings.DestinationLocales {
		destinations[i] = string(locale)
	}
	persisted := &stateV2{
		Version: 2,
		Locales: stateV2Locales{
			Source:       string(state.Settings.SourceLocale),
			Destinations: destinations,
		},
		LastUpdate: time.Now(),
		Files:      make([]stateV2File, 0, len(state.Snapshot.Documents)),
	}
	for name, document := range state.Snapshot.Documents {
		if document.Name != name {
			return nil, fmt.Errorf("snapshot document %q has mismatched name %q", name, document.Name)
		}
		digests := make(map[string]string, len(document.Digests))
		for segment, digest := range document.Digests {
			digests[segment] = string(digest)
		}
		persisted.Files = append(persisted.Files, stateV2File{Path: name, Segments: digests})
	}
	return persisted, nil
}
