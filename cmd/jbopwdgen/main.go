package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	lojban_password_gen "github.com/arran4/go-lojban-password-gen"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	gismuPath := flag.String("gismu", "", "Path to gismu.txt file")
	cmavoPath := flag.String("cmavo", "", "Path to cmavo.txt file")
	noEmbed := flag.Bool("no-embed", false, "Do not use embedded dictionaries")
	sentenceMinSize := flag.Int("minsize", 5, "Minimum number of words in the generated sentence")
	mode := flag.String("mode", "sentence", "Generation mode: 'sentence' or 'lujvo'")
	includeLujvo := flag.Bool("lujvo", false, "Include lujvo in sentence generation (only for sentence mode)")
	flag.Parse()

	findDict := func(flagPath, filename string) (io.ReadCloser, error) {
		if flagPath != "" {
			return os.Open(flagPath)
		}

		gismuEmbed, cmavoEmbed, ok := lojban_password_gen.GetEmbeddedDicts()
		if !*noEmbed && ok {
			if filename == "gismu.txt" {
				return io.NopCloser(strings.NewReader(gismuEmbed)), nil
			} else if filename == "cmavo.txt" {
				return io.NopCloser(strings.NewReader(cmavoEmbed)), nil
			}
		}

		if dictDir := os.Getenv("DICTIONARY_DIR"); dictDir != "" {
			if f, err := os.Open(dictDir + string(os.PathSeparator) + filename); err == nil {
				return f, nil
			}
		}

		if f, err := os.Open(filename); err == nil {
			return f, nil
		}

		candidates := []string{
			"/usr/share/dict/" + filename,
			"/var/share/dicts/" + filename,
		}
		for _, c := range candidates {
			if f, err := os.Open(c); err == nil {
				return f, nil
			}
		}
		return nil, fmt.Errorf("dictionary %s not found", filename)
	}

	// Parse gismu
	rGismu, err := findDict(*gismuPath, "gismu.txt")
	if err != nil {
		fmt.Println("Error finding gismu.txt:", err)
		return
	}
	defer rGismu.Close()
	gismuList, err := lojban_password_gen.ParseGismuFromReader(rGismu)
	if err != nil {
		fmt.Println("Error parsing gismu:", err)
		return
	}

	// Parse cmavo
	rCmavo, err := findDict(*cmavoPath, "cmavo.txt")
	if err != nil {
		fmt.Println("Error finding cmavo.txt:", err)
		return
	}
	defer rCmavo.Close()
	cmavoList, err := lojban_password_gen.ParseCmavoFromReader(rCmavo)
	if err != nil {
		fmt.Println("Error parsing cmavo:", err)
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
