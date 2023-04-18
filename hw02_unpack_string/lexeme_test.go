package hw02unpackstring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexeme(t *testing.T) {
	cases := []struct {
		_rune  rune
		_count uint
	}{
		{
			_rune:  'a',
			_count: 0,
		},
		{
			_rune:  'a',
			_count: 1,
		},
		{
			_rune:  'a',
			_count: 9,
		},
	}
	for _, testCase := range cases {
		noPanicFunction := func() {
			var lexeme Lexeme
			lexeme.SetRune(testCase._rune)
			lexeme.SetCount(testCase._count)

			msg := fmt.Sprintf("Lexeme %q*%d is valid. It's OK.", lexeme.GetRune(), lexeme.GetCount())
			fmt.Println(msg)
		}
		require.NotPanics(t, noPanicFunction)
	}
}

func TestLexemePanic(t *testing.T) {
	cases := []struct {
		_rune  rune
		_count uint
	}{
		{
			_rune:  'a',
			_count: 10,
		},
		{
			_rune:  'B',
			_count: 1,
		},
	}
	for _, testCase := range cases {
		panicFunction := func() {
			msg := fmt.Sprintf("Lexeme %q*%d is not valid. It's OK.", testCase._rune, testCase._count)
			fmt.Println(msg)

			var lexeme Lexeme
			lexeme.SetRune(testCase._rune)
			lexeme.SetCount(testCase._count)
		}
		require.Panics(t, panicFunction)
	}
}
func TestLexemeUnpack(t *testing.T) {
	cases := []struct {
		_rune     rune
		_count    uint
		_expected string
	}{
		{
			_rune:     'a',
			_count:    0,
			_expected: "",
		},
		{
			_rune:     'a',
			_count:    1,
			_expected: "a",
		},
		{
			_rune:     '\t',
			_count:    9,
			_expected: "\t\t\t\t\t\t\t\t\t",
		},
	}
	for _, testCase := range cases {
		testCase := testCase
		var lexeme Lexeme
		lexeme.SetRune(testCase._rune)
		lexeme.SetCount(testCase._count)
		unpacked := lexeme.Unpack()

		msg := fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", lexeme.GetRune(), lexeme.GetCount(), lexeme.Unpack())
		fmt.Println(msg)

		if unpacked != testCase._expected {
			t.Errorf(
				"(Lexeme %q*%d).Unpack() = %q; expected %q.",
				lexeme.GetRune(), lexeme.GetCount(), unpacked, testCase._expected)
		}
	}
}
