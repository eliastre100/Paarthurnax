package translation

import (
	"errors"
	"maps"
	"regexp"
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

func NewDocument(name string) *Document {
	return &Document{
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

func (document *Document) DocumentNameForLocale(locale Locale) string {
	if document.Multilingual {
		return document.Name
	}

	name := document.Name
	for _, source := range document.Locales() {

		pattern := regexp.MustCompile(`(^|[/.])` + regexp.QuoteMeta(source.String()) + `($|[/.])`)
		for { // Need to loop for overlapping occurrences (i.e.: fr/fr.yml)
			replaced := pattern.ReplaceAllString(name, "${1}"+locale.String()+"${2}")
			if replaced == name {
				break
			}
			name = replaced
		}
	}

	return name
}

func (document *Document) GetSegment(locale Locale, key string) (Segment, error) {
	catalog, ok := document.Catalogs[locale]
	if !ok {
		return Segment{}, ErrSegmentNotFound
	}
	return catalog.GetSegment(key)
}

func (document *Document) SetSegment(locale Locale, segment Segment) error {
	catalog := document.Catalog(locale)
	if catalog == nil {
		catalog = NewCatalog(locale)
		if err := document.AddCatalog(catalog); err != nil {
			return err
		}

	}
	catalog.AddSegment(segment)
	return nil
}

func (document *Document) RemoveSegment(key string) {
	for _, catalog := range document.Catalogs {
		catalog.RemoveSegment(key)
		if catalog.isEmpty() {
			document.RemoveCatalog(catalog.Locale)
		}
	}
}

func (document *Document) RemoveCatalog(locale Locale) {
	delete(document.Catalogs, locale)
}
