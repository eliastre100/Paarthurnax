package v1

import (
	v0 "Paarthurnax/internal/state/v0"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
	"os"
	"time"
)

const (
	StateFile = ".paarthurnax"
)

type State struct {
	Version int `toml:"version"`
	*Snapshot
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
		Version:  1,
		Snapshot: snapshot,
	}, nil
}

func Build(path string, srcLocale string) (*State, error) {
	s := &State{
		Version: 1,
	}
	snapshot, err := NewSnapshot(path, srcLocale)
	if err != nil {
		return nil, err
	}
	s.Snapshot = snapshot

	return s, nil
}

func (s *State) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	encoded, err := toml.Marshal(s)
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
