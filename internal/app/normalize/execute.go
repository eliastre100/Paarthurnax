package normalize

import "fmt"

func Execute(loader ProjectLoader, store DocumentStore) error {
	project, err := loader.Load()
	if err != nil {
		return fmt.Errorf("loading project: %w", err)
	}

	for _, document := range project.Documents {
		if err := store.Store(document); err != nil {
			return fmt.Errorf("storing document: %w", err)
		}
	}

	return nil
}
