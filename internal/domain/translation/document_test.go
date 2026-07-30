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
