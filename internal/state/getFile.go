package state

import "Paarthurnax/internal/state/translationFile"

func (state *State) GetFile(path string) *translationFile.TranslationFile {
	for _, file := range state.Files {
		if file.Path == path {
			return &file
		}
	}
	return nil
}
