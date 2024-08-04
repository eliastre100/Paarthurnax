package state

type State struct {
	Files []TranslationFile
}

type TranslationFile struct {
	Path           string
	SegmentsHashes map[string]string // Key: sha1
}
