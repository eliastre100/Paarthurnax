package cmd

import (
	"Paarthurnax/internal/state"
	"Paarthurnax/internal/translationGroup"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var TranslateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate the repository",
	Long:  `Translate all the modified source segment into every other language using DeepL`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Loading previous state from disk...")
		pState, err := state.Load()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Generating current state from disk...")
		nState, err := state.LoadFromDisk("config/locales", false)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Reconciling states...")
		for _, nFile := range nState.Files {
			log.Println(fmt.Sprintf("Processing %s", nFile.Path))
			changes := nFile.Changes(pState.GetFile(nFile.Path))
			if len(changes) != 0 {
				log.Println(fmt.Sprintf("Applying changes and translating %s...", nFile.Path))
				group, err := translationGroup.New(nFile.Path)
				if err != nil {
					log.Fatal(err)
				}
				if err = group.Apply(changes); err != nil {
					log.Fatal(err)
				}
			}
		}

		log.Println("Cleaning up removed files...")
		for _, pFile := range pState.Files {
			if nState.GetFile(pFile.Path) == nil {
				log.Println("Cleaning up translation of", pFile.Path)
				if errors := translationGroup.Cleanup(pFile.Path); len(errors) != 0 {
					log.Println("Unable some translation files:")
					for _, err := range errors {
						log.Println(err)
					}
				}
			}
		}

		log.Println("Persisting new state...")
		if err := nState.Save(); err != nil {
			log.Fatal("Failed to persist new state: " + err.Error())
		}
	},
}
