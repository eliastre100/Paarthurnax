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

func localesInProject(project *translation.Project) []translation.Locale {
	locales := make(map[translation.Locale]struct{})

	for _, document := range project.Documents {
		for _, locale := range document.Locales() {
			locales[locale] = struct{}{}
		}
	}

	return slices.Collect(maps.Keys(locales))
}

func selectSourceLocale(project *translation.Project, selector Selector) (translation.Locale, error) {
	projectLocales := localesInProject(project)
	if len(projectLocales) == 1 {
		return projectLocales[0], nil
	}

	choices := make(map[string]translation.Locale)
	for _, locale := range projectLocales {
		choices[locale.Name()] = locale
	}

	selectedChoice, err := selector.Select("Select the source locale", slices.Collect(maps.Keys(choices)))
	if err != nil {
		return "", err
	}
	selectedLocale, ok := choices[selectedChoice]
	if !ok {
		return "", fmt.Errorf("%s is not a valid locale", selectedChoice)
	}

	return selectedLocale, nil
}
