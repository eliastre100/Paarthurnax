package cmd

import (
	"Paarthurnax/internal/adapters/outbound/deepl"
	"Paarthurnax/internal/adapters/outbound/document"
	"Paarthurnax/internal/adapters/outbound/project"
	"Paarthurnax/internal/adapters/outbound/reporter"
	"Paarthurnax/internal/adapters/outbound/tomlproject"
	"Paarthurnax/internal/app/translate"
	deeplclient "Paarthurnax/pkg/deepl"
	"os"

	"charm.land/log/v2"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

var TranslateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate the repository",
	Long:  `Translate all the modified source segment into every other language using DeepL`,
	Run: func(cmd *cobra.Command, args []string) {
		deeplApiKey := os.Getenv("DEEPL_API_KEY")
		loader := project.NewLoader(".")
		documentStore := document.NewStore(".")
		projectStateRepository := tomlproject.NewRepository(".paarthurnax")
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

		err = translate.Execute(loader, documentStore, projectStateRepository, engine, progressReporter)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
