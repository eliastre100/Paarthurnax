package translation

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

func (translation *TranslationFile) Save() error {
	f, err := os.OpenFile(translation.Path, os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		return errors.New("Unable to open translation file: " + err.Error())
	}
	defer f.Close()

	yamlData := make(map[string]interface{})
	yamlData[translation.Locale] = translation.Segments

	yamlEncoder := yaml.NewEncoder(f)
	yamlEncoder.SetIndent(2)
	defer yamlEncoder.Close()

	if err = yamlEncoder.Encode(yamlData); err != nil {
		return errors.New("Unable to write translation file: " + err.Error())
	}
	return nil
}
