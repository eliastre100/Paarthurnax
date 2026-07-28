package tomlproject

import (
	"fmt"
	"os"

	"Paarthurnax/internal/app/settings"

	"github.com/pelletier/go-toml/v2"
)

const currentVersion = 1

var versionLoaders = map[int]func([]byte) (*settings.ProjectState, error){
	0: loadV0,
	1: loadV1,
}

type Repository struct {
	Path string
}

func NewRepository(path string) *Repository {
	return &Repository{Path: path}
}

func (r *Repository) Save(state *settings.ProjectState) error {
	persisted, err := encodeV1(state)
	if err != nil {
		return err
	}

	data, err := toml.Marshal(persisted)
	if err != nil {
		return fmt.Errorf("failed to encode project state: %w", err)
	}
	if err := os.WriteFile(r.Path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write project state: %w", err)
	}
	return nil
}

func (r *Repository) Load() (*settings.ProjectState, error) {
	data, err := os.ReadFile(r.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read project state: %w", err)
	}

	version, err := readVersion(data)
	if err != nil {
		return nil, err
	}
	loader, ok := versionLoaders[version]
	if !ok {
		return nil, fmt.Errorf("unsupported project state version: %d", version)
	}
	return loader(data)
}

func readVersion(data []byte) (int, error) {
	var header struct {
		Version int `toml:"version"`
	}
	if err := toml.Unmarshal(data, &header); err != nil {
		return 0, fmt.Errorf("failed to decode project state version: %w", err)
	}
	return header.Version, nil
}
