package locale_management

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/app/translate"
	"Paarthurnax/internal/domain/translation"
)

type fakeProjectStateRepository struct {
	state   *settings.ProjectState
	loadErr error
	saveErr error
	saved   []*settings.ProjectState
}

func (r *fakeProjectStateRepository) Load() (*settings.ProjectState, error) {
	return r.state, r.loadErr
}

func (r *fakeProjectStateRepository) Save(state *settings.ProjectState) error {
	r.saved = append(r.saved, state)
	return r.saveErr
}

type fakeSelector struct {
	selected string
	err      error
	question string
	choices  []string
	calls    int
}

func (s *fakeSelector) Select(question string, choices []string) (string, error) {
	s.calls++
	s.question = question
	s.choices = append([]string(nil), choices...)
	return s.selected, s.err
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

type translationCall struct {
	segment translation.Segment
	source  translation.Locale
	target  translation.Locale
	values  []translate.SegmentValue
}

type fakeEngine struct {
	calls  []translationCall
	result func(translation.Segment, translation.Locale) (translation.Segment, error)
}

func (e *fakeEngine) Translate(segment translation.Segment, source translation.Locale, target translation.Locale, values ...translate.SegmentValue) (translation.Segment, error) {
	e.calls = append(e.calls, translationCall{segment: segment, source: source, target: target, values: values})
	return e.result(segment, target)
}

type fakeDocumentStore struct {
	stored []*translation.Document
	err    error
}

func (s *fakeDocumentStore) Store(document *translation.Document) error {
	s.stored = append(s.stored, document)
	return s.err
}

func (*fakeDocumentStore) Delete(string) error { return nil }

type fakeReporter struct {
	changes []uint
}

func (r *fakeReporter) StartHandlingChanges(changes uint) { r.changes = append(r.changes, changes) }
func (*fakeReporter) StartHandlingDocument(string)        {}
func (*fakeReporter) UpdateDocumentChanges(string, uint)  {}
func (*fakeReporter) DoneUpdatingSegment(string, string)  {}
func (*fakeReporter) DoneUpdatingDocument(string)         {}
func (*fakeReporter) DoneDeletingDocument(string)         {}
func (*fakeReporter) DoneDeletingSegment(string, string)  {}
