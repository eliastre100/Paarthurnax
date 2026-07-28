package translation

type Project struct {
	Documents map[Locale]map[string]*Document
}

func NewProject() *Project {
	return &Project{Documents: make(map[Locale]map[string]*Document)}
}

func (p *Project) AddDocument(document *Document) {
	for _, locale := range document.Locales() {
		if _, ok := p.Documents[locale]; !ok {
			p.Documents[locale] = make(map[string]*Document)
		}
		p.Documents[locale][document.Name] = document
	}
}
