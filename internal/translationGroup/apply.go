package translationGroup

import (
	"Paarthurnax/internal/state/translationFile"
	"Paarthurnax/internal/utils"
	"Paarthurnax/pkg/deepl"
	"errors"
	"os"
	"strings"
)

func (group *TranslationGroup) Apply(changes []translationFile.Change) error {
	translator, err := deepl.New(os.Getenv("DEEPL_API_KEY"), "api-free.deepl.com")
	if err != nil {
		return errors.New("Unable to initialize the translation engine: " + err.Error())
	}
	for _, change := range changes {
		keyParts := strings.Split(change.Path, ".")
		if utils.Includes(keyParts[len(keyParts)-1], []string{"zero", "one", "other"}) {
			if err := handlePluralSegment(keyParts[len(keyParts)-1], group, &change, translator); err != nil {
				return err
			}
		} else {
			if err := handleStandaloneSegment(group, &change, translator); err != nil {
				return err
			}
		}
	}

	for _, file := range group.files {
		if err := file.Save(); err != nil {
			return errors.New("Failed to save updated translation file " + file.Path + ": " + err.Error())
		}
	}

	return nil
}

func handleStandaloneSegment(group *TranslationGroup, change *translationFile.Change, translator *deepl.Deepl) error {
	if change.Kind == translationFile.Added || change.Kind == translationFile.Updated {
		value, err := group.source.GetSegmentValueAt(change.Path)
		if err != nil {
			return errors.New("Unable to get the value of the source segment " + change.Path + ": " + err.Error())
		}
		for _, file := range group.files {
			translation, err := translator.Translate(value, "fr", file.Locale)
			if err != nil {
				return errors.New("Unable to translate the value of the source segment " + change.Path + ": " + err.Error())
			}

			if err = file.SetSegmentValueAt(change.Path, translation); err != nil {
				return errors.New("Unable to set the value of the source segment " + change.Path + " in locale " + file.Locale + ": " + err.Error())
			}
		}
	} else if change.Kind == translationFile.Removed {
		// TODO: remove key
	}
	return nil
}

func handlePluralSegment(part string, group *TranslationGroup, change *translationFile.Change, translator *deepl.Deepl) error {
	if change.Kind == translationFile.Added || change.Kind == translationFile.Updated {
		value, err := group.source.GetSegmentValueAt(change.Path)
		if err != nil {
			return errors.New("Unable to get the value of the source segment " + change.Path + ": " + err.Error())
		}

		for _, file := range group.files {
			affectedKeys := determineAffectedKeysIn(part, file.Locale)
			for _, key := range affectedKeys {
				// for each key transform the count with the language count then convert back to count
				translation, err := translator.Translate(value, "fr", file.Locale)
				if err != nil {
					return errors.New("Unable to translate the value of the source segment " + change.Path + ": " + err.Error())
				}

				pathParts := strings.Split(change.Path, ".")
				localKey := strings.Join(append(pathParts[:len(pathParts)-1], key), ".")
				if err = file.SetSegmentValueAt(localKey, translation); err != nil {
					return errors.New("Unable to set the value of the source segment " + change.Path + " in locale " + file.Locale + ": " + err.Error())
				}
			}
		}
	} else if change.Kind == translationFile.Removed {
		// TODO: remove key
	}
	return nil
}

func determineAffectedKeysIn(part string, locale string) []string {
	keysDirections := map[string]map[string][]string{
		"zero": map[string][]string{
			"en": []string{"zero"},
			"es": []string{"zero"},
			"it": []string{"zero"},
			"de": []string{"zero"},
			"hu": []string{"zero"},
			"pt": []string{"zero"},
			"pl": []string{"zero"},
			"ro": []string{"zero"},
			"uk": []string{"zero"},
		},
		"one": map[string][]string{
			"en": []string{"one"},
			"es": []string{"one"},
			"it": []string{"one"},
			"de": []string{"one"},
			"hu": []string{"one"},
			"pt": []string{"one"},
			"pl": []string{"one"},
			"ro": []string{"one"},
			"uk": []string{"one"},
		},
		"other": map[string][]string{
			"en": []string{"other"},
			"es": []string{"other"},
			"it": []string{"other"},
			"de": []string{"other"},
			"hu": []string{"other"},
			"pt": []string{"other"},
			"pl": []string{"few", "many", "other"}, // https://github.com/svenfuchs/rails-i18n/blob/master/rails/pluralization/pl.rb
			"ro": []string{"few", "other"},         // https://github.com/svenfuchs/rails-i18n/blob/master/lib/rails_i18n/common_pluralizations/romanian.rb#L5
			"uk": []string{"few", "many", "other"}, // https://github.com/svenfuchs/rails-i18n/blob/master/lib/rails_i18n/common_pluralizations/east_slavic.rb#L8
		},
	}

	return keysDirections[part][locale]
}
