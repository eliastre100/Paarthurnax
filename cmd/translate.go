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
		reporter := reporter.NewPacmanReporter()

		err = translate.Execute(loader, documentStore, projectStateRepository, engine, reporter)
		if err != nil {
			log.Fatalf("%v", err)
		}
		/*log.Info("Loading previous state from disk...", "path", v1.StateFile)
		pState, err := state.Load(v1.StateFile)
		if err != nil {
			log.Fatalf("Failed to load previous state: %v", err)
		}

		log.Info("Generating current state from disk...", "path", "config/locales")
		nState, err := state.Generate("config/locales", "fr")
		if err != nil {
			log.Fatalf("Failed to generate current state: %v", err)
		}

		log.Info("Reconciling states...")
		for _, nFile := range nState.Files {
			log.Info(fmt.Sprintf("Processing %s", nFile.Path))
			changes := nState.Compare(nFile.Path, pState.Snapshot)
			if len(changes) != 0 {
				log.Info(fmt.Sprintf("Applying changes and translating %s...", nFile.Path))
				group, err := translationgroup.NewGroup(nFile.Path)
				if err != nil {
					log.Fatal(err)
				}
				if err = group.Apply(changes); err != nil {
					log.Fatal(err)
				}
			}
		}

		log.Info("Cleaning up removed files...")
		for _, pFile := range pState.Files {
			if nState.GetFile(pFile.Path) == nil {
				log.Info("Cleaning up translation of", pFile.Path)
				if errors := translationgroup.Cleanup(pFile.Path); len(errors) != 0 {
					log.Info("Unable some translation files:")
					for _, err := range errors {
						log.Info(err)
					}
				}
			}
		}

		pState.Snapshot = nState.Snapshot
		log.Info("Persisting new state...")
		if err := pState.Save(v1.StateFile); err != nil {
			log.Fatal("Failed to persist new state: " + err.Error())
		}*/
	},
}
