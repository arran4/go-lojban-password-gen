package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	lojban_password_gen "github.com/arran4/go-lojban-password-gen"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// Define flags for configuration
	dictionaryDir := os.Getenv("DICTIONARY_DIR")
	if dictionaryDir == "" {
		dictionaryDir = "." // Default to current directory
	}

	gismuPath := flag.String("gismu", dictionaryDir+"/gismu.txt", "Path to gismu.txt file")
	cmavoPath := flag.String("cmavo", dictionaryDir+"/cmavo.txt", "Path to cmavo.txt file")
	sentenceMinSize := flag.Int("minsize", 5, "Minimum number of words in the generated sentence")
	mode := flag.String("mode", "sentence", "Generation mode: 'sentence' or 'lujvo'")
	includeLujvo := flag.Bool("lujvo", false, "Include lujvo in sentence generation (only for sentence mode)")
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

	generator := lojban_password_gen.Generator{
		GismuList:    gismuList,
		CmavoList:    cmavoList,
		IncludeLujvo: *includeLujvo,
	}

	if *mode == "lujvo" {
		lujvo, meaning := generator.GenerateLujvo()
		fmt.Printf("Generated Lujvo: %s\n", lujvo)
		fmt.Printf("Meaning: %s\n", meaning)
	} else {
		generator.GenerateSentence(*sentenceMinSize)
	}
}
