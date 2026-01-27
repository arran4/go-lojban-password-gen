package lojban_password_gen

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParseGismuFile(t *testing.T) {
	// Construct a line that matches the format expected by ParseGismuFile
	// word: 1-6 (5 chars)
	// meaning: 62-157 (95 chars)

	// 0: space
	// 1-5: "gismu"
	// 6-61: filler
	// 62+: meaning

	line := " " + fmt.Sprintf("%-5s", "gismu") + // 1-5
		strings.Repeat(" ", 56) + // 6-61 (61-6+1 = 56 chars)
		"meaning with x1 and x2 (cf. valsi)                                             "

	// Ensure line is long enough (>= 157 chars)
	if len(line) < 157 {
		line += strings.Repeat(" ", 157-len(line))
	}

	// Add keyword and hint to make it more realistic
	// Keyword: 20-40 (20 chars)
	// Hint: 41-61 (20 chars)
	// My previous string repeat overwrote these areas with spaces, which is fine as they are trimmed.
	// But let's inject them for testing.

	runes := []rune(line)
	copy(runes[20:], []rune(fmt.Sprintf("%-20s", "mykeyword")))
	copy(runes[41:], []rune(fmt.Sprintf("%-20s", "myhint")))
	line = string(runes)

	content := "Header line\n" + line + "\n"

	tmpfile, err := os.CreateTemp("", "gismu.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	list, err := ParseGismuFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseGismuFile failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Expected 1 gismu, got %d", len(list))
	}

	g := list[0]
	if g.Word != "gismu" {
		t.Errorf("Expected word 'gismu', got '%s'", g.Word)
	}
	if g.Keyword != "mykeyword" {
		t.Errorf("Expected keyword 'mykeyword', got '%s'", g.Keyword)
	}
	if g.Hint != "myhint" {
		t.Errorf("Expected hint 'myhint', got '%s'", g.Hint)
	}

	// Regex `x\d.*?` matches x1, x2
	// But `.*?` is non-greedy match of anything.
	// In "meaning with x1 and x2", `x1` matches "x1". `x2` matches "x2".
	if len(g.Placement) < 2 {
		t.Errorf("Expected at least 2 placements, got %d: %v", len(g.Placement), g.Placement)
	}

	if len(g.SeeAlso) != 1 || g.SeeAlso[0] != "valsi" {
		t.Errorf("Expected seeAlso ['valsi'], got %v", g.SeeAlso)
	}
}

func TestParseCmavoFile(t *testing.T) {
	// Construct a line for ParseCmavoFile
	// 0-11: word
	// 12-20: category
	// 21-62: keyword
	// 63+: meaning

	line := fmt.Sprintf("%-11s", "cmavo") + " " +
		fmt.Sprintf("%-8s", "cat") + " " +
		fmt.Sprintf("%-41s", "keyword") + " " +
		"meaning (cf. valsi)"

	content := "Header line\n" + line + "\n"

	tmpfile, err := os.CreateTemp("", "cmavo.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	list, err := ParseCmavoFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseCmavoFile failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Expected 1 cmavo, got %d", len(list))
	}

	c := list[0]
	if c.Word != "cmavo" {
		t.Errorf("Expected word 'cmavo', got '%s'", c.Word)
	}
	if c.Category != "cat" {
		t.Errorf("Expected category 'cat', got '%s'", c.Category)
	}
	if c.Keyword != "keyword" {
		t.Errorf("Expected keyword 'keyword', got '%s'", c.Keyword)
	}
	if len(c.SeeAlso) != 1 || c.SeeAlso[0] != "valsi" {
		t.Errorf("Expected seeAlso ['valsi'], got %v", c.SeeAlso)
	}
}

func TestGenerateSentence(t *testing.T) {
	gismuList := []Gismu{
		{Word: "gismu", Meaning: "root word"},
		{Word: "broda", Meaning: "predicate var 1"},
		{Word: "prami", Meaning: "love"},
		{Word: "ta'e", Meaning: "habitually"}, // contains '
	}
	cmavoList := []Cmavo{
		{Word: "mi", Meaning: "I"},
		{Word: "do", Meaning: "you"},
		{Word: ".ui", Meaning: "happiness"}, // contains .
		{Word: "la'o", Meaning: "the quote"}, // contains '
	}

	gen := NewGenerator(gismuList, cmavoList)

	// Test case 1: Standard generation
	sentence, _ := gen.GenerateSentence(5, false, false)
	if sentence == "" {
		t.Error("Generated sentence is empty")
	}

	// Test case 2: Include dot
	sentence, _ = gen.GenerateSentence(5, true, false)
	if !strings.HasSuffix(sentence, ".") {
		t.Errorf("Expected sentence to end with '.', got: %s", sentence)
	}

	// Test case 3: Include apostrophe
	// We run it multiple times to ensure it works consistently, as it involves randomness
	for i := 0; i < 10; i++ {
		sentence, _ = gen.GenerateSentence(5, false, true)
		if !strings.Contains(sentence, "'") {
			t.Errorf("Expected sentence to contain \"'\", got: %s", sentence)
		}
	}

	// Test case 4: Both
	sentence, _ = gen.GenerateSentence(5, true, true)
	if !strings.HasSuffix(sentence, ".") {
		t.Errorf("Expected sentence to end with '.', got: %s", sentence)
	}
	if !strings.Contains(sentence, "'") {
		t.Errorf("Expected sentence to contain \"'\", got: %s", sentence)
	}
}
