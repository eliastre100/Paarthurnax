package tomlproject

import (
	"reflect"
	"testing"

	"Paarthurnax/internal/domain/translation"
)

func TestLoadV0(t *testing.T) {
	state := loadState(t, `
[[Files]]
Path = 'config/locales/fr.yml'

[Files.SegmentsHashes]
hello = 'old-digest'
`)

	if state.Settings.SourceLocale != translation.French {
		t.Errorf("SourceLocale = %q, want %q", state.Settings.SourceLocale, translation.French)
	}
	wantDestinations := []translation.Locale{"es", "en", "de", "it", "hu", "uk", "pl", "pt", "ro"}
	if !reflect.DeepEqual(state.Settings.DestinationLocales, wantDestinations) {
		t.Errorf("DestinationLocales = %q, want %q", state.Settings.DestinationLocales, wantDestinations)
	}
	assertDocumentDigest(t, state.Snapshot, "config/locales/fr.yml", "hello", "old-digest")
}
