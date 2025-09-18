package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "how are you",
			expected: []string{"how", "are", "you"},
		},
		{
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "    I     love     VIDEO games   ",
			expected: []string{"i", "love", "video", "games"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf(`
			Faild, unmatched lenghts, got len: %v, expected len: %v
			======================================
			input: %v`, len(actual), len(c.expected), c.input)

			continue
		}

		for i := 0; i < len(actual); i++ {
			if actual[i] != c.expected[i] {
				t.Errorf(`
			Faild, unmatched world, got: %v, expected word: %v
			======================================
			input: %v
			actual: %v
			expected: %v`, actual[i], c.expected[i], c.input, actual, c.expected)

				continue
			}
		}
	}
}
