package v1

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"charm.land/log/v2"
)

type Snapshot struct {
	Date  time.Time `toml:"last_update" comment:"Date of the last successful snapshot"`
	Files []*File   `toml:"files"`
}

func NewSnapshot(path string, srcLocale string) (*Snapshot, error) {
	s := &Snapshot{
		Date:  time.Now(),
		Files: []*File{},
	}

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.Name() == srcLocale+".yml" {
			log.Debug("Processing locale", "path", path)

			f, err := NewFile(path)
			if err != nil {
				return fmt.Errorf("failed to load %s: %w", path, err)
			}

			s.Files = append(s.Files, f)
		}
		return err
	})
	return s, err
}

func (s *Snapshot) Compare(path string, other *Snapshot) []Change {
	file := s.GetFile(path)
	otherFile := other.GetFile(path)

	return file.Compare(otherFile)
}

// TODO: this is fine for now to iterate over all the files as a typical rails app will not have thousands of locale files,
// however a hash would make sense at runtime
func (s *Snapshot) GetFile(path string) *File {
	for _, file := range s.Files {
		if file.Path == path {
			return file
		}
	}
	return nil
}
