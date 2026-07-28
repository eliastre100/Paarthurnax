package translate

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"fmt"
)

func Execute(loader ProjectLoader, projectStateRepository settings.ProjectStateRepository, reporter Reporter) error {
	project, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	state, err := projectStateRepository.Load()
	if err != nil {
		return fmt.Errorf("failed to get project state: %w", err)
	}

	snapshot, err := sync.NewSnapshotFromProject(project, state.Settings.SourceLocale)
	if err != nil {
		return fmt.Errorf("failed to create project snapshot: %w", err)
	}

	for _, change := range sync.Diff(state.Snapshot, snapshot) {
		fmt.Printf("%+v\n", change)
	}

	state.Snapshot = snapshot
	return projectStateRepository.Save(state)
}
