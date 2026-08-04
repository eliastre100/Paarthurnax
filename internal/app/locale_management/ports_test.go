package locale_management

import "Paarthurnax/internal/app/settings"

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
