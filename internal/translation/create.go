package translation

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
)

func Create(path string) (*TranslationFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	r, err := regexp.Compile("(?P<locale>[a-z]+).yml")
	if err != nil {
		return nil, errors.New("Unable to prepare new translation file preprocessing: " + err.Error())
	}

	fmt.Printf("%s\n", filepath.Base(path))
	matches := r.FindStringSubmatch(filepath.Base(path))
	if matches[1] == "" {
		return nil, errors.New(path + " is an invalid translation file path")
	}
	yamlData := make(map[string]interface{})
	yamlData[matches[1]] = map[string]interface{}{}

	data, err := yaml.Marshal(yamlData)
	if err != nil {
		return nil, errors.New("Unable to create new translation file: " + err.Error())
	}
	if _, err = f.Write(data); err != nil {
		return nil, errors.New("Unable to write new translation file: " + err.Error())
	}
	if err = f.Close(); err != nil {
		return nil, errors.New("Unable to close translation file properly: " + err.Error())
	}
	return Load(path)
}
