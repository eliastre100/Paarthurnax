package translation

type Project struct {
	Documents map[string]*Document
}

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
