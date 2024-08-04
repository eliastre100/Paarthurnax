package state

import (
	"errors"
	"github.com/pelletier/go-toml/v2"
	"io"
	"os"
)

func Load() (*State, error) {
	var state State
	f, err := os.Open(".paarthurnax")
	defer f.Close()
	if err != nil {
		return nil, errors.New("Unable to open state file: " + err.Error())
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New("Unable to read state file: " + err.Error())
	}
	err = toml.Unmarshal(data, &state)
	if err != nil {
		return nil, errors.New("Unable to parse state file: " + err.Error())
	}
	return &state, nil
}
