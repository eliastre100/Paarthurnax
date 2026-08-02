package translation

import "errors"

type Project struct {
	Documents map[string]*Document
}

var ErrDocumentNotFound = errors.New("document not found")

func NewProject() *Project {
	return &Project{Documents: make(map[string]*Document)}
}

func (p *Project) AddDocument(document *Document) {
	p.Documents[document.Name] = document
}

func (p *Project) DocumentsWithCatalog(locale Locale) map[string]*Document {
	result := make(map[string]*Document)

	for _, document := range p.Documents {
		if document.HasLocale(locale) {
			result[document.Name] = document
		}
	}

	return result
}

func (p *Project) GetDocumentIn(locale Locale, name string) (*Document, error) {
	document, ok := p.Documents[name]
	if !ok {
		return nil, ErrDocumentNotFound
	}

	if document.HasLocale(locale) {
		return document, nil
	}

	destName := document.DocumentNameForLocale(locale)

	dest, ok := p.Documents[destName]
	if !ok {
		return nil, ErrDocumentNotFound
	}

	return dest, nil
}
