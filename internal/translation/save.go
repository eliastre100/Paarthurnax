package translation

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

func (translation *TranslationFile) Save() error {
	f, err := os.OpenFile(translation.Path, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return errors.New("Unable to open translation file: " + err.Error())
	}

	yamlData := make(map[string]interface{})
	yamlData[translation.Locale] = translation.Segments

	data, err := yaml.Marshal(yamlData)
	if err != nil {
		return errors.New("Unable to serialize translation file: " + err.Error())
	}
	if _, err = f.Write(data); err != nil {
		return errors.New("Unable to write translation file: " + err.Error())
	}
	if err = f.Close(); err != nil {
		return errors.New("Unable to close translation file properly: " + err.Error())
	}
	return nil
}
