package translation

type TranslationFile struct {
	Locale   string
	Path     string
	Segments map[string]interface{}
}
