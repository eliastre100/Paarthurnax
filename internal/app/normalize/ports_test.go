package normalize

import "Paarthurnax/internal/domain/translation"

type stubLoader struct {
	project *translation.Project
	err     error
}

func (l *stubLoader) Load() (*translation.Project, error) {
	return l.project, l.err
}

type recordingStore struct {
	documents []*translation.Document
	err       error
}

func (s *recordingStore) Store(document *translation.Document) error {
	s.documents = append(s.documents, document)
	return s.err
}

func (*recordingStore) Delete(string) error { return nil }
