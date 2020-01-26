package spellchecker

import (
	"regexp"
	"strings"
)

type Checker struct {
	words map[string]int
}

const (
	letters = "abcdefghijklmnopqrstuvwxyz"
)

func NewChecker(dict string) (*Checker, error) {
	words := make(map[string]int)
	checker := &Checker{words: words}

	checker.insertWords(dict)

	return checker, nil
}

func (c *Checker) insertWords(dict string) {
	re := regexp.MustCompile(`\w+`)
	for _, word := range re.FindAllString(dict, -1) {
		c.words[strings.ToLower(word)]++
	}
}

func (c *Checker) WordsCount() int {
	return len(c.words)
}

func (c *Checker) Exists(word string) bool {
	if _, ok := c.words[word]; ok {
		return true
	}

	return false
}

func (c *Checker) Corrections(word string) []string {
	if c.Exists(word) {
		return []string{word}
	}

	edits := getEdits([]string{word})
	return c.knownWords(edits)
}

func (c *Checker) knownWords(wordLists ...[]string) (known []string) {
	inserted := make(map[string]bool)

	for _, words := range wordLists {
		for _, w := range words {
			if !inserted[w] && c.Exists(w) {
				known = append(known, w)
				inserted[w] = true
			}
		}
	}

	return
}

func getSplits(word string) [][2]string {
	var splits [][2]string

	for i := 0; i < len(word)+1; i++ {
		splits = append(splits, [2]string{word[:i], word[i:]})
	}

	return splits
}

func getEdits(words []string) (edits []string) {
	for _, word := range words {
		splits := getSplits(word)
		for _, pair := range splits {
			// Delete
			if pair[1] != "" {
				edits = append(edits, pair[0]+pair[1][1:])
			}

			// Transpose
			if len(pair[1]) > 1 {
				trans := string([]byte{pair[1][1], pair[1][0]})
				edits = append(edits, pair[0]+trans+pair[1][2:])
			}

			// Replacements
			if pair[1] != "" {
				for _, c := range letters {
					edits = append(edits, pair[0]+string(c)+pair[1][1:])
				}
			}

			// Insertions
			for _, c := range letters {
				edits = append(edits, pair[0]+string(c)+pair[1])
			}
		}
	}

	return
}
