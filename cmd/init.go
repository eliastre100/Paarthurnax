package cmd

import (
	"Paarthurnax/internal/state"
	"Paarthurnax/internal/state/v1"
	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the repository",
	Long: `Initialize the repository for future use of Paarthurnax.
The repository should be in a translated state as all segments will be considered translated`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Generating state from disk...", "path", "config/locales")

		s, err := state.Generate("config/locales", "fr")
		if err != nil {
			log.Error(err)
		}

		log.Info("Persisting state...", "path", v1.StateFile)
		if err = s.Save(v1.StateFile); err != nil {
			log.Fatal(err)
		}
	},
}
