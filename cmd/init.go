package cmd

import (
	"Paarthurnax/internal/state"
	"github.com/spf13/cobra"
	"log"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the repository",
	Long: `Initialize the repository for future use of Paarthurnax.
The repository should be in a translated state as all segments will be considered translated`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating state from disk...")

		state, err := state.LoadFromDisk("config/locales", true)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Persisting state...")
		if err = state.Save(); err != nil {
			log.Fatal(err)
		}
	},
}
