package translationFile

const (
	Added   = 0
	Removed = 1
	Updated = 2
)

type Change struct {
	Kind int8
	Path string
}

func (file *TranslationFile) Changes(other *TranslationFile) []Change {
	var changes []Change

	for path, hash := range file.SegmentsHashes {
		if other == nil { // When the other translation state does not exist (new source file)
			changes = append(changes, Change{Kind: Added, Path: path})
		} else {
			oHash, ok := other.SegmentsHashes[path]
			if !ok {
				changes = append(changes, Change{Kind: Added, Path: path})
			} else if hash != oHash {
				changes = append(changes, Change{Kind: Updated, Path: path})
			}
		}
	}

	if other != nil {
		for path, _ := range other.SegmentsHashes {
			if _, ok := file.SegmentsHashes[path]; !ok {
				changes = append(changes, Change{Kind: Removed, Path: path})
			}
		}
	}

	return changes
}
