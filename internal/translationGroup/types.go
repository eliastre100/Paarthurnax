package translationGroup

import (
	"Paarthurnax/internal/translation"
)

var DestLocales = [...]string{"en"} //, "es", "it", "de"}

type TranslationGroup struct {
	Path   string
	source *translation.TranslationFile
	files  []*translation.TranslationFile
}
