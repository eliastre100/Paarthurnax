package translationFile

type TranslationFile struct {
	Path           string
	SegmentsHashes map[string]string // Key: sha1
}
