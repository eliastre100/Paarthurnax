package locales

import (
	"Paarthurnax/internal/adapters/outbound/deepl"
	"Paarthurnax/internal/adapters/outbound/document"
	"Paarthurnax/internal/adapters/outbound/project"
	"Paarthurnax/internal/adapters/outbound/reporter"
	"Paarthurnax/internal/adapters/outbound/selection"
	"Paarthurnax/internal/adapters/outbound/tomlproject"
	"Paarthurnax/internal/app/locale_management"
	"Paarthurnax/internal/app/translate"
	deeplclient "Paarthurnax/pkg/deepl"
	"os"

	"charm.land/log/v2"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [locale]",
	Short: "Add a new destination locale to the project",
	Long:  `Add a new destination locale to the project and run a translation for all existing translations`,
	Run: func(cmd *cobra.Command, args []string) {
		deeplApiKey := os.Getenv("DEEPL_API_KEY")
		if deeplApiKey == "" {
			log.Fatalf("DEEPL_API_KEY environment variable not set")
		}
		projectStateRepository := tomlproject.NewRepository(".paarthurnax")
		projectLoader := project.NewLoader(".")
		documentStore := document.NewStore(".")
		selector := selection.NewSelector()

		deeplClient, err := deeplclient.NewClient(deeplApiKey, deeplclient.DeepLDomainFree)
		if err != nil {
			log.Fatalf("Failed to create DeepL client: %v", err)
		}
		engine := deepl.NewTranslator(deeplClient)

		progressReporter := translate.Reporter(reporter.NewTextReporter())
		if term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd()) {
			pacmanReporter := reporter.NewPacmanReporter()
			progressReporter = pacmanReporter
			defer pacmanReporter.Close()
		}

		if err := locale_management.AddLocale(projectStateRepository, projectLoader, engine, progressReporter, documentStore, selector); err != nil {
			log.Fatal(err)
		}
	},
}
