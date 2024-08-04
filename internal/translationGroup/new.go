package translationGroup

import (
	"Paarthurnax/internal/translation"
	"errors"
	"strings"
)

func New(sourceFilePath string) (*TranslationGroup, error) {
	group := TranslationGroup{Path: sourceFilePath}

	source, err := translation.Load(sourceFilePath)
	if err != nil {
		return nil, errors.New("Unable to load base translation form " + sourceFilePath + ": " + err.Error())
	}
	group.source = source

	for _, locale := range DestLocales {
		file, err := translation.LoadOrCreate(strings.Replace(sourceFilePath, "fr.yml", locale+".yml", 1))
		if err != nil {
			return nil, errors.New("Unable to load " + locale + " translation form " + sourceFilePath + ": " + err.Error())
		}
		group.files = append(group.files, file)
	}

	return &group, nil
}
