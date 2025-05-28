package translation

import (
	"Paarthurnax/internal/utils"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	Locale   string
	Path     string
	Segments map[string]interface{}
}

func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read translation file: %w", err)
	}

	var yamlData map[string]interface{}
	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %w", err)
	}

	if len(yamlData) != 1 {
		return nil, fmt.Errorf("provided YAML file is not a valid translation file")
	}

	locale := utils.MapKeys(yamlData)[0]
	return &File{Path: path, Locale: locale, Segments: yamlData[locale].(map[string]interface{})}, nil
}

func Create(path string) (*File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	r, err := regexp.Compile("(?P<locale>[a-z]+).yml")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare new translation file preprocessing: %w", err)
	}

	matches := r.FindStringSubmatch(filepath.Base(path))
	if matches[1] == "" {
		return nil, fmt.Errorf("%s is an invalid translation file path", path)
	}
	yamlData := make(map[string]interface{})
	yamlData[matches[1]] = map[string]interface{}{}

	data, err := yaml.Marshal(yamlData)
	if err != nil {
		return nil, fmt.Errorf("unable to create new translation file: %w", err)
	}
	if _, err = f.Write(data); err != nil {
		return nil, fmt.Errorf("unable to write new translation file: %w", err)
	}
	if err = f.Close(); err != nil {
		return nil, fmt.Errorf("unable to close translation file properly: %w", err)
	}
	return Load(path)
}

func LoadOrCreate(filename string) (*File, error) {
	file, err := Load(filename)
	if err == nil {
		return file, nil
	}
	if strings.HasPrefix(err.Error(), "error unmarshalling YAML:") {
		return nil, err
	}
	return Create(filename)
}

func (f *File) Save() error {
	file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("unable to open translation file: %w", err)
	}
	defer file.Close()

	yamlData := make(map[string]interface{})
	yamlData[f.Locale] = f.Segments

	yamlEncoder := yaml.NewEncoder(file)
	yamlEncoder.SetIndent(2)
	defer yamlEncoder.Close()

	if err = yamlEncoder.Encode(yamlData); err != nil {
		return fmt.Errorf("unable to write translation file: %w", err)
	}
	return nil
}

func (f *File) FlattenedSegments() map[string]string {
	return extractFlattenedSegments(f.Segments)
}

func (f *File) SetSegmentValueAt(path string, value string) error {
	return setValueAt(f.Segments, path, value)
}

func (f *File) RemoveSegmentAt(path string) error {
	return removeSegmentAt(f.Segments, path)
}

func (f *File) GetSegmentValueAt(path string) (string, error) {
	return getValueAt(f.Segments, path)
}

func extractFlattenedSegments(m map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for key, value := range m {
		switch value.(type) {
		case string:
			result[key] = value.(string)
		case map[string]interface{}:
			for subKey, subValue := range extractFlattenedSegments(value.(map[string]interface{})) {
				result[key+"."+subKey] = subValue
			}
		}
	}

	return result
}

func getValueAt(m map[string]interface{}, path string) (string, error) {
	pathSegments := strings.Split(path, ".")
	if len(pathSegments) == 1 {
		value, ok := m[pathSegments[0]].(string)

		if !ok {
			return "", errors.New("path does not contain segment")
		}
		return value, nil
	}

	newM, ok := m[pathSegments[0]].(map[string]interface{})
	if !ok {
		return "", errors.New("path does not contain segment")
	}
	return getValueAt(newM, strings.Join(pathSegments[1:], "."))
}

func removeSegmentAt(m map[string]interface{}, path string) error {
	pathSegments := strings.Split(path, ".")
	if len(pathSegments) != 1 { // Trunk
		newM, ok := m[pathSegments[0]].(map[string]interface{})
		if !ok {
			return nil // The key does not exist already
		}

		if err := removeSegmentAt(newM, strings.Join(pathSegments[1:], ".")); err != nil {
			return err
		}

		if len(newM) == 0 {
			delete(m, pathSegments[0])
		}
	} else if len(pathSegments) == 1 { // Leaf
		delete(m, pathSegments[0])
	}
	return nil
}

func setValueAt(m map[string]interface{}, path string, value string) error {
	pathSegments := strings.Split(path, ".")
	if len(pathSegments) == 1 {
		m[path] = value
		return nil
	}

	newM, ok := m[pathSegments[0]].(map[string]interface{})
	if !ok {
		m[pathSegments[0]] = make(map[string]interface{})
		newM = m[pathSegments[0]].(map[string]interface{})
	}
	return setValueAt(newM, strings.Join(pathSegments[1:], "."), value)
}
