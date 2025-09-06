package state

import (
	v1 "Paarthurnax/internal/state/v1"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
)

type State interface {
	Save(path string) error
}

func Load(path string) (*v1.State, error) {
	log.Debug("Loading state form disk", "path", path)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open state file: %w", err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Error(fmt.Sprintf("Failed to close paarthurnax state file %s: %v", path, err))
		}
	}(f)

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read state file: %w", err)
	}
	version, err := getVersion(data)
	if err != nil {
		return nil, fmt.Errorf("unable to determine state version: %w", err)
	}
	log.Debug("State version determined", "path", path, "version", version)

	switch version {
	case 0:
		log.Fatal("TODO: State version 0 upgrade to v1 is not supported")
		return nil, nil
	case 1:
		return v1.Load(data)
	default:
		return nil, fmt.Errorf("state version %d is not supported", version)
	}
}

func getVersion(data []byte) (int, error) {
	var payload struct {
		Version int `toml:"version"`
	}
	if err := toml.Unmarshal(data, &payload); err != nil {
		return 0, err
	}
	return payload.Version, nil
}

func Generate(path string, srcLocale string) (*v1.State, error) {
	return v1.Build(path, srcLocale)
}
