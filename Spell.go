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

	splits := getSplits(word)

	deletes := getDeletes(word, splits)
	transposes := getTransposes(word, splits)
	replaces := getReplaces(word, splits)
	inserts := getInserts(word, splits)

	return c.knownWords(deletes, transposes, replaces, inserts)
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

func getDeletes(word string, splits [][2]string) (deletes []string) {
	for _, pair := range splits {
		if pair[1] != "" {
			deletes = append(deletes, pair[0]+pair[1][1:])
		}
	}

	return
}

func getTransposes(word string, splits [][2]string) (transposes []string) {
	for _, pair := range splits {
		if len(pair[1]) > 1 {
			trans := string([]byte{pair[1][1], pair[1][0]})
			transposes = append(transposes, pair[0]+trans+pair[1][2:])
		}
	}

	return
}

func getReplaces(word string, splits [][2]string) (replaces []string) {
	for _, pair := range splits {
		if pair[1] != "" {
			for _, c := range letters {
				replaces = append(replaces, pair[0]+string(c)+pair[1][1:])
			}
		}
	}

	return
}

func getInserts(word string, splits [][2]string) (inserts []string) {
	for _, pair := range splits {
		for _, c := range letters {
			inserts = append(inserts, pair[0]+string(c)+pair[1])
		}
	}

	return
}
