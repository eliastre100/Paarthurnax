package v0

type TranslationFile struct {
	Path           string
	SegmentsHashes map[string]string // Key: sha1
}
