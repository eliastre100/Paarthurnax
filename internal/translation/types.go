package translation

type TranslationFile struct {
	locale   string
	path     string
	segments map[string]interface{}
}
