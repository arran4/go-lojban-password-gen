package lojban_password_gen

import (
	"bufio"
	"fmt"
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

// Lists to hold parsed data
var gismuList []Gismu
var cmavoList []Cmavo

// ParseGismuFile Function to parse gismu.txt file with strict format and validation
func ParseGismuFile(filename string) ([]Gismu, error) {
	var list []Gismu
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	placementRegex := regexp.MustCompile(`x\d.*?`)         // Regex to find x1, x2...
	seeAlsoRegex := regexp.MustCompile(`\(cf\. ([^)]+)\)`) // Extract see-also references

	scanner := bufio.NewScanner(file)
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
	var list []Cmavo
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	seeAlsoRegex := regexp.MustCompile(`\(cf\. ([^)]+)\)`) // Extract see-also references
	scanner := bufio.NewScanner(file)
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

// RandomElement Function to get a random element from a slice
func RandomElement[T any](list []T) T {
	return list[rand.Intn(len(list))]
}

// GenerateSentence Generate a valid lojban sentence following grammar rules
func GenerateSentence(minSize int, includeDot bool, includeApostrophe bool) (string, []string) {
	length := rand.Intn(3) + minSize // Ensure minimum size + random expansion
	meaningDescriptions := []string{}

	var sentenceParts []string
	numberPos := rand.Intn(length + 1)

	// Add cmavo, sumti, and numbers for sentence variety
	for i := 0; i <= length; i++ {
		if numberPos == i {
			sentenceParts = append(sentenceParts, fmt.Sprint(rand.Intn(100)))
		}
		if i >= length {
			break
		}
		switch rand.Intn(10) {
		case 0, 1, 2, 3, 4:
			randSumti := RandomElement(gismuList)
			sentenceParts = append(sentenceParts, randSumti.Word)
			meaningDescriptions = append(meaningDescriptions, fmt.Sprintf("%s: %s", randSumti.Word, randSumti.Meaning))
		case 5, 6, 7, 8, 9:
			randCmavo := RandomElement(cmavoList)
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
			// Find a replacement word with an apostrophe
			var candidates []string
			var candidateMeanings []string

			// Gather candidates from gismu and cmavo
			for _, g := range gismuList {
				if strings.Contains(g.Word, "'") {
					candidates = append(candidates, g.Word)
					candidateMeanings = append(candidateMeanings, fmt.Sprintf("%s: %s", g.Word, g.Meaning))
				}
			}
			for _, c := range cmavoList {
				if strings.Contains(c.Word, "'") {
					candidates = append(candidates, c.Word)
					candidateMeanings = append(candidateMeanings, fmt.Sprintf("%s: %s", c.Word, c.Meaning))
				}
			}

			if len(candidates) > 0 {
				idx := rand.Intn(len(candidates))
				replacement := candidates[idx]
				replacementMeaning := candidateMeanings[idx]

				// Replace a random part (excluding the number if possible, or just replace valid index)
				// We need to be careful not to replace the number if we can avoid it, but numbers are parts too.
				// Numbers don't have meaning descriptions in the same list order necessarily (inserted at numberPos).
				// sentenceParts contains the words + number. meaningDescriptions contains just word meanings.

				replaceIdx := rand.Intn(len(sentenceParts))
				// Ensure we replace a word if possible, not the number
				// number is at numberPos in loop, but sentenceParts order depends on numberPos.
				// Actually sentenceParts is built sequentially.
				// if numberPos == i, we append number.
				// The number is exactly at index `numberPos` in `sentenceParts`?
				// Loop runs 0 to length.
				// if numberPos==0, number is first.
				// It seems numberPos is the index in the resulting sentenceParts.

				if replaceIdx == numberPos && len(sentenceParts) > 1 {
					// try another one
					replaceIdx = (replaceIdx + 1) % len(sentenceParts)
				}

				sentenceParts[replaceIdx] = replacement

				// Update meaning descriptions. This is trickier because number is in sentenceParts but not in meaningDescriptions.
				// We need to map sentencePart index to meaningDescription index.
				// If index < numberPos, meaningIndex = index.
				// If index > numberPos, meaningIndex = index - 1.
				// If index == numberPos, it's the number (no meaning description).

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
