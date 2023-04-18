package hw02unpackstring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexeme(t *testing.T) {
	cases := []struct {
		Rune  rune
		Count uint
	}{
		{
			Rune:  'a',
			Count: 0,
		},
		{
			Rune:  'a',
			Count: 1,
		},
		{
			Rune:  'a',
			Count: 9,
		},
	}
	for _, testCase := range cases {
		noPanicFunction := func() {
			var lexeme Lexeme
			lexeme.SetRune(testCase.Rune)
			lexeme.SetCount(testCase.Count)

			msg := fmt.Sprintf("Lexeme %q*%d is valid. It's OK.", lexeme.GetRune(), lexeme.GetCount())
			fmt.Println(msg)

		}
		require.NotPanics(t, noPanicFunction)
	}
}

func TestLexemePanic(t *testing.T) {
	cases := []struct {
		Rune  rune
		Count uint
	}{
		{
			Rune:  'a',
			Count: 10,
		},
		{
			Rune:  'B',
			Count: 1,
		},
	}
	for _, testCase := range cases {
		panicFunction := func() {
			msg := fmt.Sprintf("Lexeme %q*%d is not valid. It's OK.", testCase.Rune, testCase.Count)
			fmt.Println(msg)

			var lexeme Lexeme
			lexeme.SetRune(testCase.Rune)
			lexeme.SetCount(testCase.Count)

		}
		require.Panics(t, panicFunction)
	}
}
func TestLexemeUnpack(t *testing.T) {
	cases := []struct {
		Rune     rune
		Count    uint
		Expected string
	}{
		{
			Rune:     'a',
			Count:    0,
			Expected: "",
		},
		{
			Rune:     'a',
			Count:    1,
			Expected: "a",
		},
		{
			Rune:     '\t',
			Count:    9,
			Expected: "\t\t\t\t\t\t\t\t\t",
		},
	}
	for _, testCase := range cases {
		testCase := testCase
		var lexeme Lexeme
		lexeme.SetRune(testCase.Rune)
		lexeme.SetCount(testCase.Count)
		unpacked := lexeme.Unpack()

		msg := fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", lexeme.GetRune(), lexeme.GetCount(), lexeme.Unpack())
		fmt.Println(msg)

		if unpacked != testCase.Expected {
			t.Errorf(
				"(Lexeme %q*%d).Unpack() = %q; expected %q.",
				lexeme.GetRune(), lexeme.GetCount(), unpacked, testCase.Expected)
		}
	}
}
