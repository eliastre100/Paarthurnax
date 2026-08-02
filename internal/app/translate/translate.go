package translate

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"fmt"
)

func Execute(loader ProjectLoader, documentStore DocumentStore, projectStateRepository settings.ProjectStateRepository, engine Engine, reporter Reporter) error {
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

	changes := sync.Diff(state.Snapshot, snapshot)
	worker := NewWorker(state.Settings, project, engine, documentStore, reporter)
	reporter.StartHandlingChanges(uint(len(changes)))

	if err := process(changes, worker); err != nil {
		return fmt.Errorf("failed to process changes: %w", err)
	}

	state.Snapshot = snapshot
	return projectStateRepository.Save(state)
}

func process(changes []sync.DocumentChanges, worker *Worker) error {
	for _, documentChange := range changes {
		if err := worker.Handle(documentChange); err != nil {
			return fmt.Errorf("failed to changes from %s: %w", documentChange.Document.Name, err)
		}
	}
	return nil
}
