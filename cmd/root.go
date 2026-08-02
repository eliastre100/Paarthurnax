package cmd

import (
	"Paarthurnax/cmd/locales"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "Paarthurnax",
	Short: "Paarthurnax is a simple translation tool for Rails projects",
	Long: `A simple and quick translation tool for in place
translation of Rails project`,
	Run: func(cmd *cobra.Command, args []string) {
		println("Paarthurnax version 0.1")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			log.SetLevel(log.DebugLevel)
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(TranslateCmd)
	RootCmd.AddCommand(NormalizeCmd)
	RootCmd.AddCommand(locales.LocalesCmd)
}
