package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			input:    "   hello                                       world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    ",,,,, test, test... ",
			expected: []string{",,,,,", "test,", "test..."},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "hello\tworld\n",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(c.expected) != len(actual) {
			t.Errorf("length of actual slice: %v, length of expected slice: %v", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("%s doesnt match expected: %s", word, expectedWord)
			}
		}
	}
}
