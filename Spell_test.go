package spellchecker_test

import (
	"fmt"
	"testing"

	"spellchecker"
)

func TestSpellChecker(t *testing.T) {
	t.Run("New spell checker from words string, with correct count", func(t *testing.T) {

		cases := []struct {
			dict  string
			words []string
		}{
			{"", []string{}},
			{"the wild beast.", []string{"the", "wild", "beast"}},
			{"the The beast ..,. tHe.", []string{"the", "beast"}},
			{"don't", []string{"don", "t"}}, // XXX: Not handled as a single word
		}

		for i, test := range cases {
			name := fmt.Sprintf("Test #%d: %s", i, test.dict)
			t.Run(name, func(t *testing.T) {
				checker, err := spellchecker.NewChecker(test.dict)

				assertError(t, err, nil)
				assertWords(t, checker, test.words)
			})
		}
	})

	t.Run("Word not in checker", func(t *testing.T) {
		checker, err := spellchecker.NewChecker("")

		assertError(t, err, nil)
		if checker.WordsCount() != 0 || checker.Exists("foo") {
			t.Fatalf("expected no word in checker")
		}
	})

	t.Run("Return corrections for a word", func(t *testing.T) {
		dict := "something"
		cases := []struct {
			word       string
			correction string
		}{
			// word exits
			{"something", "something"},
			// word not found
			{"someaaang", ""},
			// one-edit
			// deletes
			{"somethingg", dict},
			// transposes
			{"osmething", dict},
			{"somehting", dict},
			{"somethign", dict},
			// replaces
			{"momething", dict},
			{"sometRing", dict},
			{"somethino", dict},
			// inserts
			{"somthing", dict},
			{"omething", dict},
			{"somethin", dict},

			// two-edits
			{"somethiaa", dict},
			{"someThin", dict},
			{"omethng", dict},
			{"somehtnig", dict},

			// three-edits, not corrected
			{"abcething", ""},
		}

		checker, err := spellchecker.NewChecker(dict)
		assertError(t, err, nil)

		for _, test := range cases {
			correction := checker.Correction(test.word)
			assertCorrection(t, correction, test.correction, test.word)
		}
	})
}

func assertError(t *testing.T, got, want error) {
	t.Helper()

	if got != want {
		t.Fatalf("Expected error %+v, got %+v", want, got)
	}
}

func assertWords(t *testing.T, checker *spellchecker.Checker, want []string) {
	t.Helper()

	count := checker.WordsCount()
	if count != len(want) {
		t.Fatalf("expected words count %d, got %d", len(want), count)
	}

	for _, word := range want {
		if !checker.Exists(word) {
			t.Fatalf("expected word %q in checker, it wasn't", word)
		}
	}
}

func assertCorrection(t *testing.T, got, want string, input string) {
	t.Helper()

	if got != want {
		t.Fatalf("Expected corrections for %q, wanted %q, got %q", input, want, got)
	}
}
