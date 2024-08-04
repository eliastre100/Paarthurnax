package translationGroup

import (
	"Paarthurnax/internal/state/translationFile"
	"Paarthurnax/pkg/deepl"
	"errors"
	"os"
)

func (group *TranslationGroup) Apply(changes []translationFile.Change) error {
	translator, err := deepl.New(os.Getenv("DEEPL_API_KEY"), "api-free.deepl.com")
	if err != nil {
		return errors.New("Unable to initialize the translation engine: " + err.Error())
	}
	for _, change := range changes {
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
	}

	for _, file := range group.files {
		if err := file.Save(); err != nil {
			return errors.New("Failed to save updated translation file " + file.Path + ": " + err.Error())
		}
	}

	return nil
}
