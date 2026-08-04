package project_initialization

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/sync"
	"Paarthurnax/internal/domain/translation"
	"fmt"
	"maps"
	"slices"

	"charm.land/log/v2"
)

func Execute(loader ProjectLoader, projectStateRepository settings.ProjectStateRepository, selector Selector) error {
	project, err := loader.Load()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	locale, err := selectSourceLocale(project, selector)
	if err != nil {
		return err
	}

	snapshot, err := sync.NewSnapshotFromProject(project, locale)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if err := projectStateRepository.Save(&settings.ProjectState{
		Settings: &settings.Settings{
			SourceLocale:       locale,
			DestinationLocales: []translation.Locale{},
		},
		Snapshot: snapshot,
	}); err != nil {
		return err
	}

	return nil
}

func localesInProject(project *translation.Project) []string {
	locales := make(map[string]struct{})

	for _, document := range project.Documents {
		for _, locale := range document.Locales() {
			locales[string(locale)] = struct{}{}
		}
	}

	return slices.Collect(maps.Keys(locales))
}

func selectSourceLocale(project *translation.Project, selector Selector) (translation.Locale, error) {
	projectLocales := localesInProject(project)
	if len(projectLocales) == 1 {
		return translation.Locale(projectLocales[0]), nil
	}

	selectedLocale, err := selector.Select("Select the source locale", projectLocales)
	if err != nil {
		return "", err
	}
	if slices.Index(projectLocales, selectedLocale) == -1 {
		return "", fmt.Errorf("%s is not a valid locale", selectedLocale)
	}

	return translation.Locale(selectedLocale), nil
}
