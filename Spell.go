package spellchecker

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Checker struct {
	words map[string]int
}

const (
	letters = "abcdefghijklmnopqrstuvwxyz"
)

func NewChecker(dict io.Reader) (*Checker, error) {
	words := make(map[string]int)
	checker := &Checker{words: words}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(dict)
	if err != nil {
		return nil, fmt.Errorf("couldn't read from dictionary: %v", err)
	}

	checker.insertWords(buf.Bytes())

	return checker, nil
}

func (c *Checker) insertWords(words []byte) {
	re := regexp.MustCompile(`\w+`)

	for _, word := range re.FindAll(words, -1) {
		c.words[strings.ToLower(string(word))]++
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

func (c *Checker) Correction(word string) string {
	if c.Exists(word) {
		return word
	}

	// one-edit
	edits := getEdits([]string{word})
	known := c.knownWords(edits)
	if len(known) > 0 {
		return c.bestCorrection(known)
	}

	// two-edits
	known = c.knownWords(getEdits(edits))
	if len(known) > 0 {
		return c.bestCorrection(known)
	}

	return ""
}

func (c *Checker) bestCorrection(corrections []string) string {
	bestCorrection, bestCount := "", 0
	for _, correction := range corrections {
		if c.words[correction] > bestCount {
			bestCount = c.words[correction]
			bestCorrection = correction
		}
	}
	return bestCorrection
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
