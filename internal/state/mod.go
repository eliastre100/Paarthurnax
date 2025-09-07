package state

import (
	v1 "Paarthurnax/internal/state/v1"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
)

const CurrentVersion = 1

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

	if version <= CurrentVersion {
		return v1.Load(data, version)
	} else {
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
