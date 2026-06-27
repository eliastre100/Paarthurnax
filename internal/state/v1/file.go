package v1

import (
	"Paarthurnax/internal/utils"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type File struct {
	Path     string            `toml:"path"`
	Segments map[string]string `toml:"segments"` // key => sha1
}

type ChangeKind int8

const (
	Added ChangeKind = iota
	Removed
	Updated
)

type Change struct {
	Kind ChangeKind
	Path string
}

func NewFile(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read translation file: %w", err)
	}

	var yamlData map[string]any
	err = yaml.Unmarshal(data, &yamlData)
	if err != nil {
		return nil, fmt.Errorf("faild to unmarshall YAML: %w", err)
	}

	if len(yamlData) != 1 {
		return nil, fmt.Errorf("multilingual locale files are not supported yet, please use single locale file")
	}

	locale := utils.MapKeys(yamlData)[0]
	return &File{
		Path:     path,
		Segments: hashTranslations(yamlData[locale].(map[string]any)),
	}, nil
}

func hashTranslations(m map[string]any) map[string]string {
	hashes := make(map[string]string)

	for k, v := range m {
		switch v.(type) {
		case string:
			h := sha1.New()
			h.Write([]byte(v.(string)))
			hashes[k] = hex.EncodeToString(h.Sum(nil))
		case map[string]interface{}:
			subHashes := hashTranslations(v.(map[string]interface{}))
			for subK, subV := range subHashes {
				hashes[k+"."+subK] = subV
			}
		}
	}

	return hashes
}

// FIXME: this should probably be on another module than the state
func (f *File) Compare(other *File) []Change {
	var changes []Change

	for path, hash := range f.Segments {
		if other == nil { // When the other translation state does not exist (new source file)
			changes = append(changes, Change{Kind: Added, Path: path})
		} else {
			oHash, ok := other.Segments[path]
			if !ok {
				changes = append(changes, Change{Kind: Added, Path: path})
			} else if hash != oHash {
				changes = append(changes, Change{Kind: Updated, Path: path})
			}
		}
	}

	if other != nil {
		for path, _ := range other.Segments {
			if _, ok := f.Segments[path]; !ok {
				changes = append(changes, Change{Kind: Removed, Path: path})
			}
		}
	}

	return changes
}
