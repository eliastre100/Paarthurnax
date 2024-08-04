package translationGroup

import (
	"Paarthurnax/internal/translation"
)

var DestLocales = [...]string{"en", "es", "it", "de", "hu", "pt", "pl", "ro", "uk"}

type TranslationGroup struct {
	Path   string
	source *translation.TranslationFile
	files  []*translation.TranslationFile
}
