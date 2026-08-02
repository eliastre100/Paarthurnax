package locales

import (
	"Paarthurnax/internal/state"
	v1 "Paarthurnax/internal/state/v1"
	"slices"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [locale]",
	Short: "Add a new destination locale to the project",
	Long:  `Add a new destination locale to the project and run a translation for all existing translations`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Error("Missing locale argument")
		}

		log.Debug("Loading state from disk...", "path", v1.StateFile)
		state, err := state.Load(v1.StateFile)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", v1.StateFile, err)
		}

		if slices.Index(state.Locales.Destinations, args[0]) != -1 {
			log.Infof("Locale %s already exists", args[0])
			return
		}

		state.Locales.Destinations = append(state.Locales.Destinations, args[0])
		if err := state.Save(v1.StateFile); err != nil {
			log.Errorf("Failed to write %s: %v", v1.StateFile, err)
		}
		log.Infof("The locale %s was added to the project added", args[0])
	},
}
