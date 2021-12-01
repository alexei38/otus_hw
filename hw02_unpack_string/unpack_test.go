package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `qw\ne`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestSiblingRunesFirstEmpty(t *testing.T) {
	const str = "Привет мир!"
	runes := []rune(str)
	prev, next := siblingRunes(runes, 0)
	require.Equal(t, "р", string(next))
	require.Equal(t, rune(0), prev)
}

func TestSiblingRunesLastEmpty(t *testing.T) {
	const str = "Привет мир!"
	runes := []rune(str)
	prev, next := siblingRunes(runes, len(runes))
	require.Equal(t, "!", string(prev))
	require.Equal(t, rune(0), next)
}

func TestSiblingRunesMid(t *testing.T) {
	const str = "Привет мир!"
	runes := []rune(str)
	prev, next := siblingRunes(runes, 3)
	require.Equal(t, "и", string(prev))
	require.Equal(t, "е", string(next))
}

func TestCheckEscapedRune(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		index    int
		expected bool
	}{
		{name: "normal_mid", input: "abcdef", index: 2, expected: false},
		{name: "normal_first", input: "abcdef", index: 0, expected: false},
		{name: "normal_last", input: "abcdef", index: 6, expected: false},
		{name: "escaped_slash1", input: `a\\bcdef`, index: 1, expected: false},
		{name: "escaped_slash2", input: `a\\bcdef`, index: 2, expected: true},
		{name: "escaped_slash3", input: `a\\bcdef`, index: 3, expected: false},
		{name: "escaped_word", input: `a\bcdef`, index: 2, expected: true},
		{name: "escaped_double1", input: `a\\\\5c\def`, index: 1, expected: false},
		{name: "escaped_double2", input: `a\\\\5c\def`, index: 2, expected: true},
		{name: "escaped_double3", input: `a\\\\5c\def`, index: 3, expected: false},
		{name: "escaped_double4", input: `a\\\\5c\def`, index: 4, expected: true},
		{name: "escaped_double5", input: `a\\\\5c\def`, index: 5, expected: false},
		{name: "escaped_double6", input: `a\\\\5c\def`, index: 8, expected: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			runes := []rune(tc.input)
			result := checkEscapedRune(runes, tc.index)
			require.Equal(t, tc.expected, result)
		})
	}
}
