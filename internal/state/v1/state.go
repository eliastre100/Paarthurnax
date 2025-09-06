package v1

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
)

const (
	StateFile = ".paarthurnax"
)

type State struct {
	Version int `toml:"version"`
	*Snapshot
}

func Load(data []byte) (*State, error) {
	var state State

	if err := toml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unable to parse state file as v1 state: %w", err)
	}
	return &state, nil
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
