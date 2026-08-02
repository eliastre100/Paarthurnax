package v1

import (
	v0 "Paarthurnax/internal/state/v0"
	"fmt"
	"os"
	"time"

	"charm.land/log/v2"
	"github.com/pelletier/go-toml/v2"
)

const (
	StateFile = ".paarthurnax"
)

type State struct {
	Version int     `toml:"version"`
	Locales Locales `toml:"locales"`
	*Snapshot
}

func New(path string, srcLocale string) (*State, error) {
	s := &State{
		Version: 1,
		Locales: Locales{
			Source:       srcLocale,
			Destinations: []string{},
		},
	}
	snapshot, err := NewSnapshot(path, srcLocale)
	if err != nil {
		return nil, err
	}
	s.Snapshot = snapshot

	return s, nil
}

func Load(data []byte, version int) (*State, error) {
	if version != 1 {
		return upgradeFromV0(data, version)
	}

	var state State

	if err := toml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unable to parse state file as v1 state: %w", err)
	}
	return &state, nil
}

func upgradeFromV0(data []byte, version int) (*State, error) {
	log.Warn("An upgrade of the state file from version 0 to version 1 is required. Upgrading now")

	v0State, err := v0.Load(data, version)
	if err != nil {
		return nil, fmt.Errorf("unable to parse state file as v0 state: %w", err)
	}

	snapshot := &Snapshot{
		Date:  time.Now(),
		Files: []*File{},
	}
	for _, f := range v0State.Files {
		snapshot.Files = append(snapshot.Files, &File{
			Path:     f.Path,
			Segments: f.SegmentsHashes,
		})
	}
	return &State{
		Version: 1,
		Locales: Locales{
			// V0 assumed that the source locale was French, and the destination locales the one mentioned
			// Given the fixed values in versions before the new state format.
			// We can safely assume that the user is using this set of locales in his project if he was using a v0 state
			Source:       "fr",
			Destinations: []string{"es", "en", "de", "it", "hu", "uk", "pl", "pt", "ro"},
		},
		Snapshot: snapshot,
	}, nil
}

func (s *State) Save(path string) error {
	encoded, err := toml.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encoded, 0644)
}
