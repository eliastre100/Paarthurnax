package translationGroup

import (
	"os"
	"strings"
)

func Cleanup(path string) []error {
	errors := make([]error, 0)

	for _, locale := range DestLocales {
		translationPath := strings.Replace(path, "fr.yml", locale+".yml", 1)
		if err := os.Remove(translationPath); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
