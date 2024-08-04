package translation

import (
	"strings"
)

func (translation *TranslationFile) SetSegmentValueAt(path string, value string) error {
	return setValueAt(translation.Segments, path, value)
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
