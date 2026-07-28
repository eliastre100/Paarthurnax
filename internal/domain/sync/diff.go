package sync

import "sort"

func Diff(previous *Snapshot, current *Snapshot) []DocumentChanges {
	previousDocuments := documents(previous)
	currentDocuments := documents(current)
	changes := make([]DocumentChanges, 0)

	for _, name := range sortedChangedDocumentNames(previousDocuments, currentDocuments) {
		previousDocument, existed := previousDocuments[name]
		currentDocument, exists := currentDocuments[name]

		switch {
		case !existed:
			changes = append(changes, documentChanges(currentDocument, ChangeTypeCreateFile))
		case !exists:
			changes = append(changes, documentChanges(previousDocument, ChangeTypeDeleteFile))
		default:
			documentChanges := diffDocument(previousDocument, currentDocument)
			if len(documentChanges) > 0 {
				changes = append(changes, DocumentChanges{
					Document: &currentDocument,
					Changes:  documentChanges,
				})
			}
		}
	}

	return changes
}

func documents(snapshot *Snapshot) map[string]Document {
	if snapshot == nil {
		return nil
	}
	return snapshot.Documents
}

func sortedChangedDocumentNames(previous, current map[string]Document) []string {
	names := make([]string, 0, len(previous)+len(current))
	for name := range previous {
		names = append(names, name)
	}
	for name := range current {
		if _, exists := previous[name]; !exists {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func documentChanges(document Document, changeType ChangeType) DocumentChanges {
	change := Change{Type: changeType, Document: &document}
	return DocumentChanges{Document: &document, Changes: []*Change{&change}}
}

func diffDocument(previous, current Document) []*Change {
	changes := make([]*Change, 0)

	for _, segment := range sortedSegmentNames(previous.Digests) {
		previousDigest := previous.Digests[segment]
		currentDigest, exists := current.Digests[segment]
		if !exists {
			changes = append(changes, segmentChange(ChangeTypeDelete, current, segment))
			continue
		}
		if previousDigest != currentDigest {
			changes = append(changes, segmentChange(ChangeTypeUpdate, current, segment))
		}
	}

	for _, segment := range sortedSegmentNames(current.Digests) {
		if _, exists := previous.Digests[segment]; !exists {
			changes = append(changes, segmentChange(ChangeTypeInsert, current, segment))
		}
	}

	return changes
}

func sortedSegmentNames(digests map[string]Digest) []string {
	segments := make([]string, 0, len(digests))
	for segment := range digests {
		segments = append(segments, segment)
	}
	sort.Strings(segments)
	return segments
}

func segmentChange(changeType ChangeType, document Document, segment string) *Change {
	return &Change{Type: changeType, Document: &document, Segment: &segment}
}
