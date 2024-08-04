package cmd

import (
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
}
