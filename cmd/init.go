package cmd

import (
	"Paarthurnax/internal/adapters/outbound/project"
	"Paarthurnax/internal/adapters/outbound/selection"
	"Paarthurnax/internal/adapters/outbound/tomlproject"
	"Paarthurnax/internal/app/project_initialization"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the repository",
	Long: `Initialize the repository for future use of Paarthurnax.
The repository should be in a translated state as all segments will be considered translated`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectStateRepository := tomlproject.NewRepository(".paarthurnax")
		projectLoader := project.NewLoader(".")
		selector := selection.NewSelector()

		err := project_initialization.Execute(projectLoader, projectStateRepository, selector)
		if err != nil {
			return err
		}
		return nil
	},
}
