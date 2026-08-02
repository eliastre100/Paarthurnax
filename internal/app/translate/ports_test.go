package translate

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/translation"
)

type translationCall struct {
	segment translation.Segment
	source  translation.Locale
	target  translation.Locale
}

type fakeEngine struct {
	calls  []translationCall
	result func(translation.Segment, translation.Locale) (translation.Segment, error)
}

func (e *fakeEngine) Translate(segment translation.Segment, source translation.Locale, target translation.Locale, _ ...SegmentValue) (translation.Segment, error) {
	e.calls = append(e.calls, translationCall{segment: segment, source: source, target: target})
	if e.result != nil {
		return e.result(segment, target)
	}
	return translation.Segment{}, nil
}

type fakeDocumentStore struct {
	stored    []*translation.Document
	deleted   []string
	storeErr  error
	deleteErr error
}

func (s *fakeDocumentStore) Store(document *translation.Document) error {
	s.stored = append(s.stored, document)
	return s.storeErr
}

func (s *fakeDocumentStore) Delete(document string) error {
	s.deleted = append(s.deleted, document)
	return s.deleteErr
}

type fakeProjectLoader struct {
	project *translation.Project
	err     error
	calls   int
}

func (l *fakeProjectLoader) Load() (*translation.Project, error) {
	l.calls++
	return l.project, l.err
}

type fakeProjectStateRepository struct {
	state *settings.ProjectState
	load  error
	save  error
	saved []*settings.ProjectState
}

func (r *fakeProjectStateRepository) Load() (*settings.ProjectState, error) {
	return r.state, r.load
}

func (r *fakeProjectStateRepository) Save(state *settings.ProjectState) error {
	r.saved = append(r.saved, state)
	return r.save
}

type fakeReporter struct {
	changes   []uint
	documents []string
	closed    bool
}

func (r *fakeReporter) StartHandlingChanges(changes uint) { r.changes = append(r.changes, changes) }
func (r *fakeReporter) StartHandlingDocument(document string) {
	r.documents = append(r.documents, document)
}
func (*fakeReporter) UpdateDocumentChanges(string, uint) {}
func (*fakeReporter) DoneUpdatingSegment(string, string) {}
func (*fakeReporter) DoneUpdatingDocument(string)        {}
func (*fakeReporter) DoneDeletingDocument(string)        {}
func (*fakeReporter) DoneDeletingSegment(string, string) {}
