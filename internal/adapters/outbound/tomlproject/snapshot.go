package tomlproject

import (
	"fmt"

	"Paarthurnax/internal/domain/sync"
)

func addDocument(snapshot *sync.Snapshot, name string, digests map[string]string) error {
	if digests == nil {
		return fmt.Errorf("snapshot document %q digests cannot be nil", name)
	}
	document := sync.NewDocument(name)
	for segment, digest := range digests {
		if err := document.AddSegmentDigest(segment, sync.Digest(digest)); err != nil {
			return fmt.Errorf("invalid snapshot document %q: %w", name, err)
		}
	}
	if err := snapshot.AddDocument(*document); err != nil {
		return fmt.Errorf("invalid project snapshot: %w", err)
	}
	return nil
}
