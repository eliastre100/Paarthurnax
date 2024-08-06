package translationGroup

import (
	"Paarthurnax/internal/state/translationFile"
	"Paarthurnax/internal/utils"
	"Paarthurnax/pkg/deepl"
	"errors"
	"fmt"
	"os"
	"regexp"
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

func prepareForTranslation(text string) (string, map[string]string, error) {
	ctx := make(map[string]string)
	r, err := regexp.Compile("%{(?P<variable>[a-zA-Z_ -]+)}")
	if err != nil {
		return "", nil, errors.New("Unable to prepare variable extraction: " + err.Error())
	}
	for _, variable := range r.FindAllStringSubmatch(text, -1) {
		ctx[variable[1]] = "<span translate=\"no\">" + variable[1] + "</span>"
		text = strings.Replace(text, "%{"+variable[1]+"}", ctx[variable[1]], -1)
	}
	return text, ctx, nil
}

func revertTranslationPreparation(text string, ctx map[string]string) (string, error) {
	for value, placeholder := range ctx {
		text = strings.Replace(text, placeholder, "%{"+value+"}", -1)
	}
	return text, nil
}

func handleStandaloneSegment(group *TranslationGroup, change *translationFile.Change, translator *deepl.Deepl) error {
	for _, file := range group.files {
		if change.Kind == translationFile.Added || change.Kind == translationFile.Updated {
			value, err := group.source.GetSegmentValueAt(change.Path)
			if err != nil {
				return errors.New("Unable to get the value of the source segment " + change.Path + ": " + err.Error())
			}
			valueForTranslation, ctx, err := prepareForTranslation(value)
			if err != nil {
				return errors.New("Unable to prepare for translation " + change.Path + ": " + err.Error())
			}
			translation, err := translator.Translate(valueForTranslation, "fr", file.Locale)
			if err != nil {
				return errors.New("Unable to translate the value of the source segment " + change.Path + ": " + err.Error())
			}
			translation, err = revertTranslationPreparation(translation, ctx)
			if err != nil {
				return errors.New("Unable to revert translation preparation for " + change.Path + ": " + err.Error())
			}

			if err = checkVariableEquity(value, translation); err != nil {
				return errors.New("The translation does not contain the required variables: " + err.Error())
			}

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
				valueForTranslation, ctx, err := prepareForTranslation(localValue)
				if err != nil {
					return errors.New("Unable to prepare for translation " + change.Path + ": " + err.Error())
				}
				translation, err := translator.Translate(valueForTranslation, "fr", file.Locale)
				if err != nil {
					return errors.New("Unable to translate the value of the source segment " + change.Path + ": " + err.Error())
				}
				translation, err = revertTranslationPreparation(translation, ctx)
				if err != nil {
					return errors.New("Unable to revert translation preparation for " + change.Path + ": " + err.Error())
				}
				translation = strings.ReplaceAll(translation, strconv.Itoa(int(definition.tip)), "%{count}")

				if err = checkVariableEquity(value, translation); err != nil {
					return errors.New("The translation does not contain the required variables: " + err.Error())
				}

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

func checkVariableEquity(a string, b string) error {
	r, err := regexp.Compile("(?P<variable>%{[a-zA-Z_ -]+})")
	if err != nil {
		return errors.New("Unable to prepare variable check: " + err.Error())
	}
	aResults := r.FindAllStringSubmatch(a, -1)

	var bResults []string
	for _, match := range r.FindAllStringSubmatch(b, -1) {
		bResults = append(bResults, match[1])
	}

	for _, variable := range aResults {
		if !utils.Includes(variable[1], bResults) {
			return errors.New(fmt.Sprintf("'%s' is not present in '%s'", variable[1], b))
		}
	}
	return nil
}
