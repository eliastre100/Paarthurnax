package tomlproject

import (
	"reflect"
	"testing"

	"Paarthurnax/internal/domain/translation"
)

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

func TestLoadV1FallsBackToLegacyLocales(t *testing.T) {
	state := loadState(t, `
version = 1

[locales]
source = ''
destinations = []

[[files]]
path = 'config/locales/fr.yml'

[files.segments]
hello = 'old-digest'
`)

	if state.Settings.SourceLocale != translation.French {
		t.Errorf("SourceLocale = %q, want %q", state.Settings.SourceLocale, translation.French)
	}
	wantDestinations := []translation.Locale{
		translation.Spanish, translation.English, translation.German, translation.Italian, translation.Hungarian, translation.Ukrainian, translation.Polish, translation.Portuguese, translation.Romanian,
	}
	if !reflect.DeepEqual(state.Settings.DestinationLocales, wantDestinations) {
		t.Errorf("DestinationLocales = %q, want %q", state.Settings.DestinationLocales, wantDestinations)
	}
	assertDocumentDigest(t, state.Snapshot, "config/locales/fr.yml", "hello", "old-digest")
}
