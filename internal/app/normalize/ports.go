package normalize

import "Paarthurnax/internal/domain/translation"

type DocumentStore interface {
	Store(document *translation.Document) error
	Delete(document string) error
}

type ProjectLoader interface {
	Load() (*translation.Project, error)
}
