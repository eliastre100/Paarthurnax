package translation

import (
	"Paarthurnax/internal/utils"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

func Load(path string) (*TranslationFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("Error reading file:" + err.Error())
	}

	var yamlData map[string]interface{}
	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return nil, errors.New("Error unmarshaling YAML:" + err.Error())
	}

	if len(yamlData) != 1 {
		return nil, errors.New("the provided YAML file is not a valid translation file")
	}

	locale := utils.MapKeys(yamlData)[0]
	file := TranslationFile{Path: path, Locale: locale, Segments: yamlData[locale].(map[string]interface{})}

	return &file, nil
}
