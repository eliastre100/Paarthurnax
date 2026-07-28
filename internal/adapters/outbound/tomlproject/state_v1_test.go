package tomlproject

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"
)

func TestV1SaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".paarthurnax")
	repository := NewRepository(path)
	state := newProjectState(t)

	if err := repository.Save(state); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	for _, key := range []string{"version = 1", "[locales]", "[[files]]", "[files.segments]"} {
		if !strings.Contains(string(contents), key) {
			t.Errorf("saved TOML = %q, want %q", contents, key)
		}
	}

	got, err := repository.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, state) {
		t.Errorf("Load() = %#v, want %#v", got, state)
	}
}

func TestV1SaveAndLoadEmptySnapshot(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".paarthurnax")
	repository := NewRepository(path)
	state := &settings.ProjectState{
		Settings: &settings.Settings{
			SourceLocale:       translation.English,
			DestinationLocales: []translation.Locale{},
		},
		Snapshot: sync.NewSnapshot("en"),
	}

	if err := repository.Save(state); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	got, err := repository.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, state) {
		t.Errorf("Load() = %#v, want %#v", got, state)
	}
}

func TestLoadV1(t *testing.T) {
	state := loadState(t, `
version = 1
last_update = 2025-09-06T18:20:28Z

[locales]
source = 'ja'
destinations = ['en', 'fr']

[[files]]
path = 'config/locales/ja.yml'

[files.segments]
hello = 'new-digest'
`)

	if state.Settings.SourceLocale != translation.Japanese {
		t.Errorf("SourceLocale = %q, want %q", state.Settings.SourceLocale, translation.Japanese)
	}
	wantDestinations := []translation.Locale{translation.English, translation.French}
	if !reflect.DeepEqual(state.Settings.DestinationLocales, wantDestinations) {
		t.Errorf("DestinationLocales = %q, want %q", state.Settings.DestinationLocales, wantDestinations)
	}
	assertDocumentDigest(t, state.Snapshot, "config/locales/ja.yml", "hello", "new-digest")
}

func TestV1SaveRejectsInvalidState(t *testing.T) {
	repository := NewRepository(filepath.Join(t.TempDir(), ".paarthurnax"))

	for name, state := range map[string]*settings.ProjectState{
		"nil state":    nil,
		"nil settings": {Snapshot: sync.NewSnapshot("en")},
		"nil snapshot": {Settings: &settings.Settings{}},
	} {
		t.Run(name, func(t *testing.T) {
			if err := repository.Save(state); err == nil {
				t.Error("Save() error = nil, want invalid state error")
			}
		})
	}
}
