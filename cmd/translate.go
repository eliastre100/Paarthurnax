package cmd

import (
	"Paarthurnax/internal/state"
	"Paarthurnax/internal/state/v1"
	"Paarthurnax/internal/translationgroup"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var TranslateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate the repository",
	Long:  `Translate all the modified source segment into every other language using DeepL`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Loading previous state from disk...", "path", v1.StateFile)
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
		}
	},
}
