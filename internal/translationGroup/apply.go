package translationGroup

import (
	"Paarthurnax/internal/state/translationFile"
	"Paarthurnax/internal/utils"
	"Paarthurnax/pkg/deepl"
	"errors"
	"os"
	"strconv"
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
	for _, file := range group.files {
		if change.Kind == translationFile.Added || change.Kind == translationFile.Updated {
			value, err := group.source.GetSegmentValueAt(change.Path)
			if err != nil {
				return errors.New("Unable to get the value of the source segment " + change.Path + ": " + err.Error())
			}
			translation, err := translator.Translate(value, "fr", file.Locale)
			if err != nil {
				return errors.New("Unable to translate the value of the source segment " + change.Path + ": " + err.Error())
			}

			// TODO: Add check variables all presents

			if err = file.SetSegmentValueAt(change.Path, translation); err != nil {
				return errors.New("Unable to set the value of the source segment " + change.Path + " in locale " + file.Locale + ": " + err.Error())
			}
		} else if change.Kind == translationFile.Removed {
			if err := file.RemoveSegmentAt(change.Path); err != nil {
				return errors.New("Unable to remove segment " + change.Path + " in locale " + file.Locale + ": " + err.Error())
			}
		}
	}
	return nil
}

func handlePluralSegment(part string, group *TranslationGroup, change *translationFile.Change, translator *deepl.Deepl) error {
	for _, file := range group.files {
		affectedKeys := determineAffectedKeysIn(part, file.Locale)
		for _, definition := range affectedKeys {
			pathParts := strings.Split(change.Path, ".")
			localKey := strings.Join(append(pathParts[:len(pathParts)-1], definition.key), ".")

			if change.Kind == translationFile.Added || change.Kind == translationFile.Updated {
				value, err := group.source.GetSegmentValueAt(change.Path)
				if err != nil {
					return errors.New("Unable to get the value of the source segment " + change.Path + ": " + err.Error())
				}

				localValue := strings.ReplaceAll(value, "%{count}", strconv.Itoa(int(definition.tip)))
				translation, err := translator.Translate(localValue, "fr", file.Locale)
				if err != nil {
					return errors.New("Unable to translate the value of the source segment " + change.Path + ": " + err.Error())
				}
				translation = strings.ReplaceAll(translation, strconv.Itoa(int(definition.tip)), "%{count}")

				// TODO: Add check variables all presents
				if err = file.SetSegmentValueAt(localKey, translation); err != nil {
					return errors.New("Unable to set the value of the source segment " + change.Path + " in locale " + file.Locale + ": " + err.Error())
				}
			} else if change.Kind == translationFile.Removed {
				if err := file.RemoveSegmentAt(localKey); err != nil {
					return errors.New("Unable to remove segment " + localKey + " in locale " + file.Locale + ": " + err.Error())
				}
			}
		}
	}
	return nil
}

type pluralDefinition struct {
	key string
	tip int32
}

func determineAffectedKeysIn(part string, locale string) []pluralDefinition {
	keysDirections := map[string]map[string][]pluralDefinition{
		"zero": {
			"en": {
				{
					key: "zero",
					tip: 0,
				},
			},
			"es": {
				{
					key: "zero",
					tip: 0,
				},
			},
			"it": {
				{
					key: "zero",
					tip: 0,
				},
			},
			"de": {
				{
					key: "zero",
					tip: 0,
				},
			},
			"hu": {
				{
					key: "zero",
					tip: 0,
				},
			},
			"pt": {
				{
					key: "zero",
					tip: 0,
				},
			},
			"pl": {
				{
					key: "zero",
					tip: 0,
				},
			},
			"ro": {
				{
					key: "zero",
					tip: 0,
				},
			},
			"uk": {
				{
					key: "zero",
					tip: 0,
				},
			},
		},
		"one": {
			"en": {
				{
					key: "one",
					tip: 1,
				},
			},
			"es": {
				{
					key: "one",
					tip: 1,
				},
			},
			"it": {
				{
					key: "one",
					tip: 1,
				},
			},
			"de": {
				{
					key: "one",
					tip: 1,
				},
			},
			"hu": {
				{
					key: "one",
					tip: 1,
				},
			},
			"pt": {
				{
					key: "one",
					tip: 1,
				},
			},
			"pl": {
				{
					key: "one",
					tip: 1,
				},
			},
			"ro": {
				{
					key: "one",
					tip: 1,
				},
			},
			"uk": {
				{
					key: "one",
					tip: 1,
				},
			},
		},
		"other": {
			"en": {
				{
					key: "other",
					tip: 2,
				},
			},
			"es": {
				{
					key: "other",
					tip: 2,
				},
			},
			"it": {
				{
					key: "other",
					tip: 2,
				},
			},
			"de": {
				{
					key: "other",
					tip: 2,
				},
			},
			"hu": {
				{
					key: "other",
					tip: 2,
				},
			},
			"pt": {
				{
					key: "other",
					tip: 2,
				},
			},
			"pl": {
				{
					key: "few",
					tip: 2,
				},
				{
					key: "many",
					tip: 11,
				},
				{
					key: "other",
					tip: 2, // No value match
				},
			}, // https://github.com/svenfuchs/rails-i18n/blob/master/rails/pluralization/pl.rb
			"ro": {
				{
					key: "few",
					tip: 2,
				},
				{
					key: "other",
					tip: 20,
				},
			}, // https://github.com/svenfuchs/rails-i18n/blob/master/lib/rails_i18n/common_pluralizations/romanian.rb#L5
			"uk": {
				{
					key: "few",
					tip: 2,
				},
				{
					key: "many",
					tip: 20,
				},
				{
					key: "other",
					tip: 11,
				},
			}, // https://github.com/svenfuchs/rails-i18n/blob/master/lib/rails_i18n/common_pluralizations/east_slavic.rb#L8
		},
	}

	return keysDirections[part][locale]
}
