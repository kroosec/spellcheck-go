package spellchecker_test

import (
	"fmt"
	"reflect"
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
			word        string
			corrections []string
		}{
			// word exits
			{"something", []string{"something"}},
			// word not found
			{"someaaang", nil},
			// deletes
			{"somethingg", []string{"something"}},
			// transposes
			{"osmething", []string{"something"}},
			{"somehting", []string{"something"}},
			{"somethign", []string{"something"}},
			// replaces
			{"momething", []string{"something"}},
			{"sometRing", []string{"something"}},
			{"somethino", []string{"something"}},
			// inserts
			{"somthing", []string{"something"}},
			{"omething", []string{"something"}},
			{"somethin", []string{"something"}},
		}
		checker, err := spellchecker.NewChecker(dict)
		assertError(t, err, nil)

		for _, test := range cases {
			corrections := checker.Corrections(test.word)
			assertCorrections(t, corrections, test.corrections, test.word)
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

func assertCorrections(t *testing.T, got, want []string, input string) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Expected corrections for %q, wanted %+v, got %+v", input, want, got)
	}
}
