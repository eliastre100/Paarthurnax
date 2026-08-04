package locales

import (
	"Paarthurnax/internal/adapters/outbound/project"
	"Paarthurnax/internal/adapters/outbound/selection"
	"Paarthurnax/internal/adapters/outbound/tomlproject"
	"Paarthurnax/internal/app/locale_management"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [locale]",
	Short: "Add a new destination locale to the project",
	Long:  `Add a new destination locale to the project and run a translation for all existing translations`,
	Run: func(cmd *cobra.Command, args []string) {
		projectStateRepository := tomlproject.NewRepository(".paarthurnax")
		projectLoader := project.NewLoader(".")
		selector := selection.NewSelector()

		if err := locale_management.AddLocale(projectStateRepository, projectLoader, selector); err != nil {
			log.Fatal(err)
		}
	},
}
