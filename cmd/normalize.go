package cmd

import (
	"Paarthurnax/internal/state"
	"Paarthurnax/internal/translation"
	"Paarthurnax/internal/translationGroup"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var NormalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize the repository",
	Long:  `Normalize all the other language to limit noise on sub-secant translations`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Loading current state from disk...")
		nState, err := state.LoadFromDisk("config/locales", false)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Normalizing translations...")
		for _, nFile := range nState.Files {
			log.Println(fmt.Sprintf("Processing %s...", nFile.Path))

			for _, locale := range translationGroup.DestLocales {
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

		log.Println("Done!")
	},
}
