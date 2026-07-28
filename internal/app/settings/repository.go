package settings

import "Paarthurnax/internal/domain/sync"

type ProjectState struct {
	Settings *Settings
	Snapshot *sync.Snapshot
}

type ProjectStateRepository interface {
	Load() (*ProjectState, error)
	Save(state *ProjectState) error
}
