package state

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

func (state *State) Save() error {
	f, err := os.Create(".paarthurnax")
	if err != nil {
		return err
	}
	encoded, err := toml.Marshal(state)
	if err != nil {
		return err
	}
	if _, err = f.Write(encoded); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}
