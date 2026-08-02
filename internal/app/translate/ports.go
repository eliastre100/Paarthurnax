package translate

import "Paarthurnax/internal/domain/translation"

type SegmentValue struct {
	Name  string
	Value string
}

type Engine interface {
	Translate(segment translation.Segment, srcLocale translation.Locale, dstLocale translation.Locale, values ...SegmentValue) (translation.Segment, error)
}

type DocumentStore interface {
	Store(document *translation.Document) error
	Delete(document string) error
}

type ProjectLoader interface {
	Load() (*translation.Project, error)
}

type Reporter interface {
	StartHandlingChanges(documentChanged uint)
	StartHandlingDocument(document string)
	UpdateDocumentChanges(document string, count uint)
	DoneUpdatingSegment(document string, name string)
	DoneUpdatingDocument(document string)
	DoneDeletingDocument(document string)
	DoneDeletingSegment(document string, name string)
}
