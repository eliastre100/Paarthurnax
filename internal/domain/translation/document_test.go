package translation

import (
	"errors"
	"reflect"
	"slices"
	"testing"
)

func TestNewDocument(t *testing.T) {
	document := NewDocument("foo")

	if document.ID == "" {
		t.Errorf("ID = %q, want a generated ID", document.ID)
	}
	if document.Name != "foo" {
		t.Errorf("Name = %q, want foo", document.Name)
	}
	if document.Multilingual {
		t.Error("Multilingual = true, want false")
	}
	if document.Catalogs == nil {
		t.Fatal("Catalogs = nil, want initialized map")
	}
	if len(document.Catalogs) != 0 {
		t.Errorf("len(Catalogs) = %d, want 0", len(document.Catalogs))
	}
}

func TestDocumentAddCatalog(t *testing.T) {
	t.Run("adds a single catalog", func(t *testing.T) {
		document := NewDocument("foo")
		catalog := Catalog{ID: "catalog-en", Locale: "en"}

		err := document.AddCatalog(&catalog)

		if err != nil {
			t.Fatalf("AddCatalog() error = %v, want nil", err)
		}
		if len(document.Catalogs) != 1 {
			t.Fatalf("len(Catalogs) = %d, want 1", len(document.Catalogs))
		}
		if got := document.Catalogs[catalog.Locale]; !reflect.DeepEqual(*got, catalog) {
			t.Errorf("Catalogs[%q] = %#v, want %#v", catalog.Locale, got, catalog)
		}
		if document.Multilingual {
			t.Error("Multilingual = true, want false")
		}
	})

	t.Run("adds two catalogs and makes the document multilingual", func(t *testing.T) {
		document := NewDocument("foo")
		catalogs := []Catalog{
			{ID: "catalog-en", Locale: "en"},
			{ID: "catalog-fr", Locale: "fr"},
		}

		for _, catalog := range catalogs {
			if err := document.AddCatalog(&catalog); err != nil {
				t.Fatalf("AddCatalog(%q) error = %v, want nil", catalog.Locale, err)
			}
		}

		if len(document.Catalogs) != 2 {
			t.Fatalf("len(Catalogs) = %d, want 2", len(document.Catalogs))
		}
		if !document.Multilingual {
			t.Error("Multilingual = false, want true")
		}
	})

	t.Run("rejects a second catalog with the same locale", func(t *testing.T) {
		document := NewDocument("foo")
		first := Catalog{ID: "catalog-en-1", Locale: "en"}
		duplicate := Catalog{ID: "catalog-en-2", Locale: "en"}
		if err := document.AddCatalog(&first); err != nil {
			t.Fatalf("AddCatalog() first error = %v, want nil", err)
		}

		err := document.AddCatalog(&duplicate)

		if !errors.Is(err, ErrCatalogAlreadyExists) {
			t.Fatalf("AddCatalog() error = %v, want %v", err, ErrCatalogAlreadyExists)
		}
		if len(document.Catalogs) != 1 {
			t.Errorf("len(Catalogs) = %d, want 1", len(document.Catalogs))
		}
		if got := document.Catalogs[first.Locale]; got.ID != first.ID {
			t.Errorf("Catalogs[%q] = %#v, want original catalog %#v", first.Locale, got, first)
		}
		if document.Multilingual {
			t.Error("Multilingual = true, want false")
		}
	})
}

func TestDocumentLocales(t *testing.T) {
	t.Run("returns no locales when it has no catalogs", func(t *testing.T) {
		document := NewDocument("foo")

		locales := document.Locales()

		if len(locales) != 0 {
			t.Errorf("Locales() = %v, want no locales", locales)
		}
	})

	t.Run("returns the locale of each catalog", func(t *testing.T) {
		document := NewDocument("foo")
		for _, locale := range []Locale{English, French, Japanese} {
			if err := document.AddCatalog(&Catalog{Locale: locale}); err != nil {
				t.Fatalf("AddCatalog(%q) error = %v, want nil", locale, err)
			}
		}

		locales := document.Locales()

		slices.Sort(locales)
		if want := []Locale{English, French, Japanese}; !slices.Equal(locales, want) {
			t.Errorf("Locales() = %v, want %v", locales, want)
		}
	})
}

func TestDocumentHasLocale(t *testing.T) {
	document := NewDocument("foo")
	english := NewCatalog(English)
	if err := document.AddCatalog(english); err != nil {
		t.Fatalf("AddCatalog() error = %v, want nil", err)
	}

	if !document.HasLocale(English) {
		t.Error("HasLocale(English) = false, want true")
	}
	if document.HasLocale(French) {
		t.Error("HasLocale(French) = true, want false")
	}
}

func TestDocumentCatalog(t *testing.T) {
	document := NewDocument("foo")
	english := NewCatalog(English)
	if err := document.AddCatalog(english); err != nil {
		t.Fatalf("AddCatalog() error = %v, want nil", err)
	}

	if got := document.Catalog(English); got != english {
		t.Errorf("Catalog(English) = %p, want %p", got, english)
	}
	if got := document.Catalog(French); got != nil {
		t.Errorf("Catalog(French) = %p, want nil", got)
	}
}

func TestDocumentGetSegment(t *testing.T) {
	document := NewDocument("messages.yml")
	segment := Segment{Key: "greeting", Value: "Hello"}
	english := NewCatalog(English)
	english.AddSegment(segment)
	if err := document.AddCatalog(english); err != nil {
		t.Fatalf("AddCatalog() error = %v, want nil", err)
	}

	got, err := document.GetSegment(English, segment.Key)
	if err != nil {
		t.Fatalf("GetSegment() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, segment) {
		t.Errorf("GetSegment() = %#v, want %#v", got, segment)
	}
}

func TestDocumentGetSegmentNotFound(t *testing.T) {
	document := NewDocument("messages.yml")

	got, err := document.GetSegment(French, "missing")

	if !errors.Is(err, ErrSegmentNotFound) {
		t.Fatalf("GetSegment() error = %v, want %v", err, ErrSegmentNotFound)
	}
	if !reflect.DeepEqual(got, Segment{}) {
		t.Errorf("GetSegment() = %#v, want zero Segment", got)
	}
}

func TestDocumentSetSegment(t *testing.T) {
	document := NewDocument("messages.yml")
	original := Segment{Key: "greeting", Value: "Hello"}
	updated := Segment{Key: "greeting", Value: "Bonjour", Plural: true}

	if err := document.SetSegment(French, original); err != nil {
		t.Fatalf("SetSegment() error = %v, want nil", err)
	}
	if !document.HasLocale(French) {
		t.Fatal("SetSegment() did not create the catalog for the requested locale")
	}

	if err := document.SetSegment(French, updated); err != nil {
		t.Fatalf("SetSegment() replacement error = %v, want nil", err)
	}

	catalog := document.Catalog(French)
	if len(catalog.Segments) != 1 {
		t.Fatalf("len(Segments) = %d, want 1", len(catalog.Segments))
	}
	if got := catalog.Segments[updated.Key]; !reflect.DeepEqual(got, updated) {
		t.Errorf("Segments[%q] = %#v, want %#v", updated.Key, got, updated)
	}
}

func TestDocumentRemoveSegment(t *testing.T) {
	document := NewDocument("messages.yml")
	for _, locale := range []Locale{English, French} {
		catalog := NewCatalog(locale)
		catalog.AddSegment(Segment{Key: "removed", Value: "remove me"})
		catalog.AddSegment(Segment{Key: "kept", Value: "keep me"})
		if err := document.AddCatalog(catalog); err != nil {
			t.Fatalf("AddCatalog(%q) error = %v, want nil", locale, err)
		}
	}

	document.RemoveSegment("removed")

	for _, locale := range []Locale{English, French} {
		catalog := document.Catalog(locale)
		if catalog == nil {
			t.Fatalf("Catalog(%q) = nil, want catalog containing remaining segment", locale)
		}
		if _, err := catalog.GetSegment("removed"); !errors.Is(err, ErrSegmentNotFound) {
			t.Errorf("GetSegment(%q, removed) error = %v, want %v", locale, err, ErrSegmentNotFound)
		}
		if _, err := catalog.GetSegment("kept"); err != nil {
			t.Errorf("GetSegment(%q, kept) error = %v, want nil", locale, err)
		}
	}
}

func TestDocumentRemoveSegmentRemovesEmptyCatalogs(t *testing.T) {
	document := NewDocument("messages.yml")
	for _, locale := range []Locale{English, French} {
		if err := document.SetSegment(locale, Segment{Key: "greeting", Value: "Hello"}); err != nil {
			t.Fatalf("SetSegment(%q) error = %v, want nil", locale, err)
		}
	}

	document.RemoveSegment("greeting")

	if len(document.Catalogs) != 0 {
		t.Errorf("len(Catalogs) = %d, want 0", len(document.Catalogs))
	}
}

func TestDocumentRemoveCatalog(t *testing.T) {
	document := NewDocument("messages.yml")
	for _, locale := range []Locale{English, French} {
		if err := document.AddCatalog(NewCatalog(locale)); err != nil {
			t.Fatalf("AddCatalog(%q) error = %v, want nil", locale, err)
		}
	}

	document.RemoveCatalog(French)

	if document.HasLocale(French) {
		t.Error("HasLocale(French) = true, want false after RemoveCatalog")
	}
	if !document.HasLocale(English) {
		t.Error("HasLocale(English) = false, want true after removing another catalog")
	}
}

func TestDocumentNameForLocale(t *testing.T) {
	for _, tt := range []struct {
		name         string
		documentName string
		multilingual bool
		sourceLocale Locale
		locale       Locale
		want         string
	}{
		{
			name:         "replaces a locale-only filename",
			documentName: "en.yml",
			sourceLocale: English,
			locale:       French,
			want:         "fr.yml",
		},
		{
			name:         "replaces a locale filename prefix",
			documentName: "en.messages.yml",
			sourceLocale: English,
			locale:       French,
			want:         "fr.messages.yml",
		},
		{
			name:         "replaces a locale filename suffix",
			documentName: "messages.en.yml",
			sourceLocale: English,
			locale:       French,
			want:         "messages.fr.yml",
		},
		{
			name:         "replaces a locale directory",
			documentName: "en/messages.yml",
			sourceLocale: English,
			locale:       French,
			want:         "fr/messages.yml",
		},
		{
			name:         "replaces every occurrence of a locale",
			documentName: "en/en.yml",
			sourceLocale: English,
			locale:       French,
			want:         "fr/fr.yml",
		},
		{
			name:         "preserves parent directories",
			documentName: "config/locales/en.messages.yml",
			sourceLocale: English,
			locale:       Japanese,
			want:         "config/locales/ja.messages.yml",
		},
		{
			name:         "preserves names without a locale pattern",
			documentName: "config/messages.yml",
			sourceLocale: English,
			locale:       French,
			want:         "config/messages.yml",
		},
		{
			name:         "does not treat a non-catalog filename token as a locale",
			documentName: "do.another.yml",
			sourceLocale: English,
			locale:       French,
			want:         "do.another.yml",
		},
		{
			name:         "replaces a locale between filename delimiters",
			documentName: "ex.am.ple.yml",
			sourceLocale: "am",
			locale:       French,
			want:         "ex.fr.ple.yml",
		},
		{
			name:         "replaces a locale between path delimiters",
			documentName: "ex/am/ple",
			sourceLocale: "am",
			locale:       French,
			want:         "ex/fr/ple",
		},
		{
			name:         "uses the catalog locale even when it is not a known locale constant",
			documentName: "do.another.yml",
			sourceLocale: "another",
			locale:       French,
			want:         "do.fr.yml",
		},
		{
			name:         "preserves names when there is no source catalog",
			documentName: "en.yml",
			locale:       French,
			want:         "en.yml",
		},
		{
			name:         "preserves multilingual document names",
			documentName: "en.messages.yml",
			multilingual: true,
			sourceLocale: English,
			locale:       French,
			want:         "en.messages.yml",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			document := Document{Name: tt.documentName, Multilingual: tt.multilingual}
			if tt.sourceLocale != "" {
				document.Catalogs = map[Locale]*Catalog{tt.sourceLocale: {Locale: tt.sourceLocale}}
			}

			if got := document.DocumentNameForLocale(tt.locale); got != tt.want {
				t.Errorf("DocumentNameForLocale(%q) = %q, want %q", tt.locale, got, tt.want)
			}
		})
	}
}
