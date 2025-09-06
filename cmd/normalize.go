package cmd

import (
	"Paarthurnax/internal/state"
	"Paarthurnax/internal/translation"
	"Paarthurnax/internal/translationgroup"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"strings"
)

var NormalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize the repository",
	Long:  `Normalize all the other language to limit noise on sub-secant translations`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Loading current state from disk...")
		nState, err := state.Generate("config/locales", "fr")
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Normalizing translations...")
		for _, nFile := range nState.Files {
			log.Info(fmt.Sprintf("Processing %s...", nFile.Path))

			for _, locale := range translationgroup.DestLocales {
				path := strings.Replace(nFile.Path, "fr.yml", locale+".yml", 1)
				file, err := translation.LoadOrCreate(path)
				if err != nil {
					log.Fatal(fmt.Sprintf("%s: %s", path, err.Error()))
				}
				if err = file.Save(); err != nil {
					log.Fatal(err)
				}
			}
		}

		log.Info("Done!")
	},
}
