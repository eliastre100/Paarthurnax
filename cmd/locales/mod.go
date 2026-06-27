package locales

import (
	"github.com/spf13/cobra"
)

var LocalesCmd = &cobra.Command{
	Use:   "locales",
	Short: "Handle locales configuration",
	Long:  `locales subcommands allow to manage the locales configuration of the project`,
}

func init() {
	LocalesCmd.AddCommand(AddCmd)
}
