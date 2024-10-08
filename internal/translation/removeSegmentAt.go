package translation

import (
	"strings"
)

func (translation *TranslationFile) RemoveSegmentAt(path string) error {
	return removeSegmentAt(translation.Segments, path)
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
