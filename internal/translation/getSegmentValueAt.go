package translation

import (
	"errors"
	"strings"
)

func (translation *TranslationFile) GetSegmentValueAt(path string) (string, error) {
	return getValueAt(translation.Segments, path)
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
