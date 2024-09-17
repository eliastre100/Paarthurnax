package translation

import (
	"strings"
)

func LoadOrCreate(filename string) (*TranslationFile, error) {
	file, err := Load(filename)
	if err == nil {
		return file, nil
	}
	if strings.HasPrefix(err.Error(), "Error unmarshaling YAML:") {
		return nil, err
	}
	return Create(filename)
}
