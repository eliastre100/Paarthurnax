package locale_management

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/app/translate"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"
	"errors"
	"fmt"
	"maps"
	"slices"
)

func AddLocale(projectStateRepository settings.ProjectStateRepository, projectLoader ProjectLoader, translator translate.Engine, reporter translate.Reporter, documentStore translate.DocumentStore, selector Selector) error {
	state, err := projectStateRepository.Load()
	if err != nil {
		return fmt.Errorf("could not load project state: %w", err)
	}

	project, err := projectLoader.Load()
	if err != nil {
		return fmt.Errorf("could not load project: %w", err)
	}

	locale, err := selectNewLocale(selector, *state.Settings)
	if err != nil {
		return err
	}
	state.Settings.DestinationLocales = append(state.Settings.DestinationLocales, locale)

	worker := translate.NewWorker(&settings.Settings{
		SourceLocale:       state.Settings.SourceLocale,
		DestinationLocales: []translation.Locale{locale},
	}, project, translator, documentStore, reporter)
	missingTranslations := buildDocumentChangesForNewLocale(project, locale, state.Snapshot)
	if err := process(reporter, missingTranslations, worker); err != nil {
		return fmt.Errorf("could not process missing translations: %w", err)
	}

	// TODO: translate segments missing in the new destination compared to the state

	if err := projectStateRepository.Save(state); err != nil {
		return fmt.Errorf("could not save project state: %w", err)
	}
	return nil
}

func selectNewLocale(selector Selector, s settings.Settings) (translation.Locale, error) {
	availableLocales := make(map[string]translation.Locale)

	for _, locale := range translation.Locales {
		availableLocales[locale.Name()] = locale
	}

	delete(availableLocales, s.SourceLocale.Name())
	for _, locale := range s.DestinationLocales {
		delete(availableLocales, locale.Name())
	}

	choice, err := selector.Select("Select the locale to add as destination", slices.Collect(maps.Keys(availableLocales)))
	if err != nil {
		return "", err
	}

	locale, ok := availableLocales[choice]
	if !ok {
		return "", fmt.Errorf("%s is not a valid choice", choice)
	}

	return locale, nil
}

func buildDocumentChangesForNewLocale(project *translation.Project, locale translation.Locale, snapshot *sync.Snapshot) []sync.DocumentChanges {
	var actions []sync.DocumentChanges

	for _, snapshotDocument := range snapshot.Documents {
		changes := sync.DocumentChanges{
			Document: &snapshotDocument,
			Changes:  make([]*sync.Change, 0),
		}

		document, err := project.GetDocumentIn(locale, snapshotDocument.Name)
		if errors.Is(err, translation.ErrDocumentNotFound) {
			changes.Changes = append(changes.Changes, &sync.Change{
				Type:     sync.ChangeTypeCreateFile,
				Document: &snapshotDocument,
				Segment:  nil,
			})
			actions = append(actions, changes)
			continue
		}

		for key, _ := range snapshotDocument.Digests {
			_, err := document.GetSegment(locale, key)
			if errors.Is(err, translation.ErrSegmentNotFound) {
				changes.Changes = append(changes.Changes, &sync.Change{
					Type:     sync.ChangeTypeInsert,
					Document: &snapshotDocument,
					Segment:  &key,
				})
			}
		}

		if len(changes.Changes) > 0 {
			actions = append(actions, changes)
		}
	}

	return actions
}

func process(reporter translate.Reporter, missingTranslations []sync.DocumentChanges, worker *translate.Worker) error {
	reporter.StartHandlingChanges(uint(len(missingTranslations)))

	for _, change := range missingTranslations {
		if err := worker.Handle(change); err != nil {
			return err
		}
	}
	return nil
}
