package locale_management

import (
	"Paarthurnax/internal/app/settings"
	"Paarthurnax/internal/domain/translation"
	"fmt"
	"maps"
	"slices"
)

func AddLocale(projectStateRepository settings.ProjectStateRepository, projectLoader ProjectLoader, selector Selector) error {
	state, err := projectStateRepository.Load()
	if err != nil {
		return fmt.Errorf("could not load project state: %w", err)
	}

	//project, err := projectLoader.Load()
	//if err != nil {
	//	return fmt.Errorf("could not load project: %w", err)
	//}

	locale, err := selectNewLocale(selector, *state.Settings)
	if err != nil {
		return err
	}
	state.Settings.DestinationLocales = append(state.Settings.DestinationLocales, locale)

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
