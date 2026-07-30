package translation

import (
	"errors"
	"maps"
	"slices"

	"github.com/google/uuid"
)

type Document struct {
	ID           string
	Name         string
	Multilingual bool
	Catalogs     map[Locale]*Catalog
}

var ErrCatalogAlreadyExists = errors.New("catalog already exists")

func NewDocument(name string) Document {
	return Document{
		ID:       uuid.New().String(),
		Name:     name,
		Catalogs: make(map[Locale]*Catalog),
	}
}

func (document *Document) AddCatalog(catalog *Catalog) error {
	if _, ok := document.Catalogs[catalog.Locale]; ok {
		return ErrCatalogAlreadyExists
	}

	document.Catalogs[catalog.Locale] = catalog
	if len(document.Catalogs) > 1 {
		document.Multilingual = true // FIXME: this is weak as a new multilingual document can be created and only have the source locale at first
	}
	return nil
}

func (document *Document) Locales() []Locale {
	return slices.Collect(maps.Keys(document.Catalogs))
}

func (document *Document) HasLocale(locale Locale) bool {
	_, ok := document.Catalogs[locale]
	return ok
}

func (document *Document) Catalog(locale Locale) *Catalog {
	return document.Catalogs[locale]
}
