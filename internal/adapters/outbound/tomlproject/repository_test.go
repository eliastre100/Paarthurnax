package tomlproject

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepositoryLoadErrors(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		_, err := NewRepository(filepath.Join(t.TempDir(), "missing")).Load()
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Load() error = %v, want wrapped os.ErrNotExist", err)
		}
	})

	t.Run("invalid TOML", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), ".paarthurnax")
		if err := os.WriteFile(path, []byte("[settings\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := NewRepository(path).Load()
		if err == nil || !strings.Contains(err.Error(), "decode") {
			t.Errorf("Load() error = %v, want TOML decoding error", err)
		}
	})

	t.Run("unsupported version", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), ".paarthurnax")
		if err := os.WriteFile(path, []byte("version = 2\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := NewRepository(path).Load()
		if err == nil || !strings.Contains(err.Error(), "unsupported project state version: 2") {
			t.Errorf("Load() error = %v, want unsupported version error", err)
		}
	})
}
