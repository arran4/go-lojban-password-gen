package lojban_password_gen

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

type Gismu struct {
	Word      string
	RafsiCVC  string
	RafsiCCV  string
	RafsiCVV  string
	Keyword   string
	Hint      string
	Meaning   string
	Placement []string // Extracted x1, x2, etc.
	SeeAlso   []string
}

type Cmavo struct {
	Word     string
	Category string
	Keyword  string
	Meaning  string
	SeeAlso  []string
}

// Generator holds the dictionary lists and caching for efficient generation
type Generator struct {
	GismuList           []Gismu
	CmavoList           []Cmavo
	GismuWithApostrophe []Gismu
	CmavoWithApostrophe []Cmavo
	IncludeLujvo bool
}

// NewGenerator initializes a new Generator with the provided lists
func NewGenerator(gismu []Gismu, cmavo []Cmavo) *Generator {
	gen := &Generator{
		GismuList: gismu,
		CmavoList: cmavo,
	}

	// Pre-calculate lists of words with apostrophes
	for _, g := range gismu {
		if strings.Contains(g.Word, "'") {
			gen.GismuWithApostrophe = append(gen.GismuWithApostrophe, g)
		}
	}
	for _, c := range cmavo {
		if strings.Contains(c.Word, "'") {
			gen.CmavoWithApostrophe = append(gen.CmavoWithApostrophe, c)
		}
	}

	return gen
}

// ParseGismuFile Function to parse gismu.txt file with strict format and validation
func ParseGismuFile(filename string) ([]Gismu, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	return ParseGismuFromReader(file)
}

func ParseGismuFromReader(r io.Reader) ([]Gismu, error) {
	var list []Gismu
	placementRegex := regexp.MustCompile(`x\d.*?`)         // Regex to find x1, x2...
	seeAlsoRegex := regexp.MustCompile(`\(cf\. ([^)]+)\)`) // Extract see-also references

	scanner := bufio.NewScanner(r)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()

		// Skip version line or header
		if lineNumber == 1 {
			continue
		}

		if len(line) < 157 {
			return nil, fmt.Errorf("line %d too short (expected >=157): %s", lineNumber, line)
		}

		word := strings.TrimSpace(line[1:6])
		rafsiCVC := strings.TrimSpace(line[7:10])
		rafsiCCV := strings.TrimSpace(line[11:14])
		rafsiCVV := strings.TrimSpace(line[15:19])
		keyword := strings.TrimSpace(line[20:40])
		hint := strings.TrimSpace(line[41:61])
		meaning := strings.TrimSpace(line[62:157])
		placements := placementRegex.FindAllString(meaning, -1)

		// Extract see-also references if present
		seeAlsoMatches := seeAlsoRegex.FindStringSubmatch(line)
		var seeAlso []string
		if len(seeAlsoMatches) > 1 {
			seeAlso = strings.Split(seeAlsoMatches[1], ", ")
		}

		list = append(list, Gismu{
			Word:      word,
			RafsiCVC:  rafsiCVC,
			RafsiCCV:  rafsiCCV,
			RafsiCVV:  rafsiCVV,
			Keyword:   keyword,
			Hint:      hint,
			Meaning:   meaning,
			Placement: placements,
			SeeAlso:   seeAlso,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return list, nil
}

// ParseCmavoFile Function to parse cmavo.txt file with validation
func ParseCmavoFile(filename string) ([]Cmavo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()
	return ParseCmavoFromReader(file)
}

func ParseCmavoFromReader(r io.Reader) ([]Cmavo, error) {
	var list []Cmavo
	seeAlsoRegex := regexp.MustCompile(`\(cf\. ([^)]+)\)`) // Extract see-also references
	scanner := bufio.NewScanner(r)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()

		// Skip version line or header
		if lineNumber == 1 {
			continue
		}

		if len(line) < 63 {
			return nil, fmt.Errorf("line %d too short (expected >=63): %s", lineNumber, line)
		}

		word := strings.TrimSpace(line[0:11])      // Columns 0-11: Word
		category := strings.TrimSpace(line[12:20]) // Columns 12-20: Category
		keyword := strings.TrimSpace(line[21:62])  // Columns 21-62: Keyword
		meaning := strings.TrimSpace(line[63:])    // Columns 63+: Meaning

		// Extract see-also references if present
		seeAlsoMatches := seeAlsoRegex.FindStringSubmatch(line)
		var seeAlso []string
		if len(seeAlsoMatches) > 1 {
			seeAlso = strings.Split(seeAlsoMatches[1], ", ")
		}

		list = append(list, Cmavo{
			Word:     word,
			Category: category,
			Keyword:  keyword,
			Meaning:  meaning,
			SeeAlso:  seeAlso,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return list, nil
}

// cryptoIntn returns a random integer in [0, max) using crypto/rand
func cryptoIntn(max int) int {
	if max <= 0 {
		return 0
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		// Fallback or panic? For a password generator, panic on failure is safer than weak randomness.
		panic(fmt.Sprintf("failed to get secure random number: %v", err))
	}
	return int(n.Int64())
}

// RandomElement Function to get a random element from a slice
func RandomElement[T any](list []T) T {
	if len(list) == 0 {
		var zero T
		return zero
	}
	return list[cryptoIntn(len(list))]
}

// GenerateLujvo Generates a logical compound word (lujvo) from 2 random gismu
func (g *Generator) GenerateLujvo() (string, string) {
	// Pick 2 gismu
	g1 := RandomElement(g.GismuList)
	g2 := RandomElement(g.GismuList)

	// Helper to get a rafsi
	getRafsi := func(gis Gismu, isLast bool) string {
		var candidates []string
		if gis.RafsiCVC != "" {
			candidates = append(candidates, gis.RafsiCVC)
		}
		if gis.RafsiCCV != "" {
			candidates = append(candidates, gis.RafsiCCV)
		}
		if gis.RafsiCVV != "" {
			candidates = append(candidates, gis.RafsiCVV)
		}

		if len(candidates) > 0 {
			return candidates[rand.Intn(len(candidates))]
		}
		// Fallback
		if isLast {
			return gis.Word // Use full word at end
		}
		// Crude fallback: 4 letters + y
		if len(gis.Word) >= 4 {
			return gis.Word[0:4] + "y"
		}
		return gis.Word // Should not happen for valid gismu
	}

	r1 := getRafsi(g1, false)
	r2 := getRafsi(g2, true)

	// Simple concatenation
	lujvo := r1 + r2
	meaning := fmt.Sprintf("lujvo(%s + %s)", g1.Keyword, g2.Keyword)
	return lujvo, meaning
}

// GenerateSentence Generate a valid lojban sentence following grammar rules
func (g *Generator) GenerateSentence(minSize int, includeDot bool, includeApostrophe bool) (string, []string) {
	length := cryptoIntn(3) + minSize // Ensure minimum size + random expansion
	meaningDescriptions := []string{}

	var sentenceParts []string
	numberPos := cryptoIntn(length + 1)

	// Add cmavo, sumti, and numbers for sentence variety
	for i := 0; i <= length; i++ {
		if numberPos == i {
			sentenceParts = append(sentenceParts, fmt.Sprint(cryptoIntn(100)))
		}
		if i >= length {
			break
		}
		r := cryptoIntn(10)
		if g.IncludeLujvo && rand.Intn(5) == 0 {
			lujvo, meaning := g.GenerateLujvo()
			sentenceParts = append(sentenceParts, lujvo)
			meaningDescriptions = append(meaningDescriptions, fmt.Sprintf("%s: %s", lujvo, meaning))
			continue
		}

		switch r {
		case 0, 1, 2, 3, 4:
			randSumti := RandomElement(g.GismuList)
			sentenceParts = append(sentenceParts, randSumti.Word)
			meaningDescriptions = append(meaningDescriptions, fmt.Sprintf("%s: %s", randSumti.Word, randSumti.Meaning))
		case 5, 6, 7, 8, 9:
			randCmavo := RandomElement(g.CmavoList)
			sentenceParts = append(sentenceParts, randCmavo.Word)
			meaningDescriptions = append(meaningDescriptions, fmt.Sprintf("%s: %s", randCmavo.Word, randCmavo.Meaning))
		}
	}

	if includeApostrophe {
		hasApostrophe := false
		for _, part := range sentenceParts {
			if strings.Contains(part, "'") {
				hasApostrophe = true
				break
			}
		}

		if !hasApostrophe && len(sentenceParts) > 0 {
			var candidates []string
			var candidateMeanings []string

			// Use pre-calculated lists
			for _, item := range g.GismuWithApostrophe {
				candidates = append(candidates, item.Word)
				candidateMeanings = append(candidateMeanings, fmt.Sprintf("%s: %s", item.Word, item.Meaning))
			}
			for _, item := range g.CmavoWithApostrophe {
				candidates = append(candidates, item.Word)
				candidateMeanings = append(candidateMeanings, fmt.Sprintf("%s: %s", item.Word, item.Meaning))
			}

			if len(candidates) > 0 {
				idx := cryptoIntn(len(candidates))
				replacement := candidates[idx]
				replacementMeaning := candidateMeanings[idx]

				replaceIdx := cryptoIntn(len(sentenceParts))
				// Ensure we replace a word if possible, not the number
				// numberPos is the index in the resulting sentenceParts.

				if replaceIdx == numberPos && len(sentenceParts) > 1 {
					// try another one
					replaceIdx = (replaceIdx + 1) % len(sentenceParts)
				}

				sentenceParts[replaceIdx] = replacement

				// Update meaning descriptions. The number is in sentenceParts at numberPos, but not in meaningDescriptions.
				if replaceIdx != numberPos {
					meaningIdx := replaceIdx
					if replaceIdx > numberPos {
						meaningIdx--
					}
					if meaningIdx >= 0 && meaningIdx < len(meaningDescriptions) {
						meaningDescriptions[meaningIdx] = replacementMeaning
					}
				}
			}
		}
	}

	sentence := strings.Join(sentenceParts, " ")
	if includeDot {
		sentence += "."
	}

	return sentence, meaningDescriptions
}
