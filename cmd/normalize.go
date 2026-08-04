package cmd

import (
	"Paarthurnax/internal/adapters/outbound/document"
	"Paarthurnax/internal/adapters/outbound/project"
	"Paarthurnax/internal/app/normalize"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

var NormalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize the repository",
	Long:  `Normalize all the other language to limit noise on sub-secant translations`,
	Run: func(cmd *cobra.Command, args []string) {
		loader := project.NewLoader(".")
		documentStore := document.NewStore(".")

		err := normalize.Execute(loader, documentStore)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
