package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"

	"strconv"
)

func contains(Runes []rune, Rune rune) bool {
	for _, r := range Runes {
		if r == Rune {
			return true
		}
	}
	return false
}

// Not valid "ABCDEFGHIJKLMNOPQRSTUVWXYZ".
var validRunes = []rune("abcdefghijklmnopqrstuvwxyz\n\t") 

func isRune(Rune rune) (rune, error) {
	if !contains(validRunes, Rune) {
		msg := fmt.Sprintf("Not valid rune: '%c'", Rune)
		return Rune, errors.New(msg)
	}
	return Rune, nil
}

func isCount(Count uint) (uint, error) {
	if Count > 9 {
		msg := fmt.Sprintf("Count '%d' must be between 0 and 9.", Count)
		return Count, errors.New(msg)
	}
	return Count, nil
}

func RuneToUint(r rune) (uint, error) {
	u, errParseUint := strconv.ParseUint(string(r), 10, 64)
	return uint(u), errParseUint
}

type Lexeme struct {
	_rune        rune
	_runeWasSet  bool
	_count       uint
	_countWasSet bool
}

func (lexeme *Lexeme) SetCount(count uint) {
	_, errCount := isCount(count)
	if errCount != nil {
		panic(errCount)
	}
	lexeme._count = count
	lexeme._countWasSet = true
}

func (lexeme *Lexeme) GetCount() uint {
	return lexeme._count
}

func (lexeme *Lexeme) IsCountWasSet() bool {
	return lexeme._countWasSet
}

func (lexeme *Lexeme) SetRune(Rune rune) {
	_, err := isRune(Rune)
	if err != nil {
		panic(err)
	}
	lexeme._rune = Rune
	lexeme._runeWasSet = true
}

func (lexeme *Lexeme) GetRune() rune {
	return lexeme._rune
}

func (lexeme *Lexeme) IsRuneWasSet() bool {
	return lexeme._runeWasSet
}

func (lexeme *Lexeme) Unpack() string {
	count := lexeme._count
	if !lexeme._countWasSet {
		count = 1
	}
	return strings.Repeat(string(lexeme._rune), int(count))
}

func (lexeme *Lexeme) isEmpty() bool {
	if lexeme._runeWasSet || lexeme._countWasSet {
		return false
	}
	return true
}
