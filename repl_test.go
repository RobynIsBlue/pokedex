package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hi welcome    to  DOGS",
			expected: []string{"hi", "welcome", "to", "dogs"},
		},
		{
			input:    "sup",
			expected: []string{"sup"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			actual := cleanInput(c.input)
			require.Equal(t, len(c.expected), len(actual))
			for i := range actual {
				word := actual[i]
				expectedWord := c.expected[i]
				assert.Equal(t, expectedWord, word)
			}
		})
	}
}

func TestAPIGet(t *testing.T) {
	cases := []struct {
		input    any
		expected locationData
	}{
		{
			input: 2,
			expected: locationData{
				ID:   2,
				Name: "eterna-city-area",
			},
		},
		{
			input: "eterna-city-area",
			expected: locationData{
				ID:   2,
				Name: "eterna-city-area",
			},
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%d", c.input), func(t *testing.T) {
			actual, err := callLocationApi(c.input)
			require.NoError(t, err)
			require.Equal(t, c.expected, actual)
		})
	}
}
