package main

import (
	"flag"
	"fmt"
	lojban_password_gen "github.com/arran4/go-lojban-password-gen"
	"os"
)

func main() {
	// Define flags for configuration
	dictionaryDir := os.Getenv("DICTIONARY_DIR")
	if dictionaryDir == "" {
		dictionaryDir = "." // Default to current directory
	}

	gismuPath := flag.String("gismu", dictionaryDir+"/gismu.txt", "Path to gismu.txt file")
	cmavoPath := flag.String("cmavo", dictionaryDir+"/cmavo.txt", "Path to cmavo.txt file")
	sentenceMinSize := flag.Int("minsize", 5, "Minimum number of words in the generated sentence")
	flag.Parse()

	// Parse gismu.txt and cmavo.txt files
	gismuList, err := lojban_password_gen.ParseGismuFile(*gismuPath)
	if err != nil {
		fmt.Println("Error parsing gismu file:", err)
		return
	}

	cmavoList, err := lojban_password_gen.ParseCmavoFile(*cmavoPath)
	if err != nil {
		fmt.Println("Error parsing cmavo file:", err)
		return
	}

	if len(gismuList) == 0 || len(cmavoList) == 0 {
		fmt.Println("No data loaded. Please check gismu.txt and cmavo.txt files.")
		return
	}

	lojban_password_gen.GenerateSentence(*sentenceMinSize)
}
