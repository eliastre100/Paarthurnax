package translationgroup

import (
	v1 "Paarthurnax/internal/state/v1"
	"Paarthurnax/internal/translation"
	"Paarthurnax/internal/utils"
	"Paarthurnax/pkg/deepl"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var DestLocales = [...]string{"en", "es", "it", "de", "hu", "pt", "pl", "ro", "uk"}

type Group struct {
	Path   string
	source *translation.File
	files  []*translation.File
}

func NewGroup(sourceFilePath string) (*Group, error) {
	group := Group{Path: sourceFilePath}

	source, err := translation.Load(sourceFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to load base translation form %s: %w", sourceFilePath, err)
	}
	group.source = source

	for _, locale := range DestLocales {
		file, err := translation.LoadOrCreate(strings.Replace(sourceFilePath, "fr.yml", locale+".yml", 1))
		if err != nil {
			return nil, fmt.Errorf("unable to load %s translation form %s: %w", locale, sourceFilePath, err)
		}
		group.files = append(group.files, file)
	}

	return &group, nil
}

func (group *Group) Apply(changes []v1.Change) error {
	translator, err := deepl.NewClient(os.Getenv("DEEPL_API_KEY"), deepl.DeepLDomainFree)
	if err != nil {
		return fmt.Errorf("unable to initialize the translation engine: %w", err)
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
			return fmt.Errorf("failed to save updated translation file %s: %w", file.Path, err)
		}
	}

	return nil
}

func prepareForTranslation(text string) (string, map[string]string, error) {
	ctx := make(map[string]string)
	r, err := regexp.Compile("%{(?P<variable>[a-zA-Z_ -]+)}")
	if err != nil {
		return "", nil, fmt.Errorf("unable to prepare variable extraction: %w", err)
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

func handleStandaloneSegment(group *Group, change *v1.Change, translator *deepl.Client) error {
	for _, file := range group.files {
		if change.Kind == v1.Added || change.Kind == v1.Updated {
			value, err := group.source.GetSegmentValueAt(change.Path)
			if err != nil {
				return fmt.Errorf("unable to get the value of the source segment %s: %w", change.Path, err)
			}
			valueForTranslation, ctx, err := prepareForTranslation(value)
			if err != nil {
				return fmt.Errorf("unable to prepare for translation %s: %w", change.Path, err)
			}
			translation, err := translator.Translate(valueForTranslation, "fr", file.Locale)
			if err != nil {
				return fmt.Errorf("unable to translate the value of the source segment %s: %w", change.Path, err.Error())
			}
			translation, err = revertTranslationPreparation(translation, ctx)
			if err != nil {
				return fmt.Errorf("unable to revert translation preparation for %s: %w", change.Path, err.Error())
			}

			if err = checkVariableEquity(value, translation); err != nil {
				return fmt.Errorf("the translation does not contain the required variables: %w", err)
			}

			if err = file.SetSegmentValueAt(change.Path, translation); err != nil {
				return fmt.Errorf("unable to set the value of the source segment %s in locale %s: %w", change.Path, file.Locale, err)
			}
		} else if change.Kind == v1.Removed {
			if err := file.RemoveSegmentAt(change.Path); err != nil {
				return fmt.Errorf("unable to remove segment %s in locale %s: %w", change.Path, file.Locale, err)
			}
		}
	}
	return nil
}

func handlePluralSegment(part string, group *Group, change *v1.Change, translator *deepl.Client) error {
	for _, file := range group.files {
		affectedKeys := determineAffectedKeysIn(part, file.Locale)
		for _, definition := range affectedKeys {
			pathParts := strings.Split(change.Path, ".")
			localKey := strings.Join(append(pathParts[:len(pathParts)-1], definition.key), ".")

			if change.Kind == v1.Added || change.Kind == v1.Updated {
				value, err := group.source.GetSegmentValueAt(change.Path)
				if err != nil {
					return fmt.Errorf("unable to get the value of the source segment %s: %w", change.Path, err)
				}

				countTip := fmt.Sprintf("<span translate=\"no\">%d</span>", definition.tip)
				localValue := strings.ReplaceAll(value, "%{count}", countTip)
				valueForTranslation, ctx, err := prepareForTranslation(localValue)
				if err != nil {
					return fmt.Errorf("unable to prepare for translation %s: %w", change.Path, err.Error())
				}
				translation, err := translator.Translate(valueForTranslation, "fr", file.Locale)
				if err != nil {
					return fmt.Errorf("unable to translate the value of the source segment %s: %w", change.Path, err.Error())
				}
				translation, err = revertTranslationPreparation(translation, ctx)
				if err != nil {
					return fmt.Errorf("unable to revert translation preparation for %s: %w", change.Path, err)
				}
				translation = strings.ReplaceAll(translation, countTip, "%{count}")

				if err = checkVariableEquity(value, translation); err != nil {
					return fmt.Errorf("the translation does not contain the required variables: %w", err)
				}

				if err = file.SetSegmentValueAt(localKey, translation); err != nil {
					return fmt.Errorf("unable to set the value of the source segment %s in locale %s: %w", change.Path, file.Locale, err)
				}
			} else if change.Kind == v1.Removed {
				if err := file.RemoveSegmentAt(localKey); err != nil {
					return fmt.Errorf("unable to remove segment %s in locale %s: %w", localKey, file.Locale, err)
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
	return keysDirections[part][locale]
}

func checkVariableEquity(a string, b string) error {
	r, err := regexp.Compile("(?P<variable>%{[a-zA-Z_ -]+})")
	if err != nil {
		return fmt.Errorf("unable to prepare variable check: %w", err)
	}
	aResults := r.FindAllStringSubmatch(a, -1)

	var bResults []string
	for _, match := range r.FindAllStringSubmatch(b, -1) {
		bResults = append(bResults, match[1])
	}

	for _, variable := range aResults {
		if !utils.Includes(variable[1], bResults) {
			return fmt.Errorf("'%s' is not present in '%s'", variable[1], b)
		}
	}
	return nil
}

func Cleanup(path string) []error {
	errors := make([]error, 0)

	for _, locale := range DestLocales {
		translationPath := strings.Replace(path, "fr.yml", locale+".yml", 1)
		if err := os.Remove(translationPath); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
