package hw02unpackstring

import (
	"errors"
	"fmt"
)

func Unpack(inputString string) (string, error) {
	length := len(inputString)

	// TODO Выгоднее лучше не делать реверс строки, а идти последовательно назад
	// for position, Rune := range stringutil.Reverse(inputString) {
	var currentBlock *StatementBlock
	lexeme := Lexeme{}
	for position, _ := range inputString {
		reverseIndex := length - position - 1
		symbol := rune(inputString[reverseIndex])
		symbolUint, errOfRuneToUint := RuneToUint(symbol)
		if errOfRuneToUint == nil {
			symbolCount, errOfCount := isCount(symbolUint)
			if errOfCount == nil {
				if lexeme.IsCountWasSet() {
					msg := fmt.Sprintf("Not valid count-symbol '%d' in position `%d`.", lexeme._count, reverseIndex+1)
					return "", errors.New(msg)
				}
				lexeme.SetCount(symbolCount)
				continue
			}
		}

		symbolRune, errOfRune := isRune(symbol)
		if errOfRune == nil {
			lexeme.SetRune(symbolRune)
			var currentBlockTmp StatementBlock = StatementBlock{BlockLexeme: lexeme, NextBlock: currentBlock}
			currentBlock = &currentBlockTmp
			lexeme = Lexeme{}
			continue
		}

		msg := fmt.Sprintf("Not valid symbol '%c' in position `%d`.", symbol, reverseIndex)
		return "", errors.New(msg)

	}

	if !lexeme.isEmpty() {
		msg := fmt.Sprintf("Not valid symbol '%c' in position `%d`.", inputString[0], 0)
		return "", errors.New(msg)
	}

	return currentBlock.Unpack(), nil
}

func main() {
	// unpacked, err := Unpack("3abc")
	// fmt.Println("err", err)
	// fmt.Println("unpacked", unpacked)
	// fmt.Println("")

	// unpacked, err := Unpack("45")
	// fmt.Println("err", err)
	// fmt.Println("unpacked", unpacked)
	// fmt.Println("")

	// unpacked, err := Unpack("aaa10b")
	// fmt.Println("err", err)
	// fmt.Println("unpacked", unpacked)
	// fmt.Println("")

	// unpacked, err := Unpack("a4bc2d5e")
	// fmt.Println("err", err)
	// fmt.Println("unpacked", unpacked)
	// fmt.Println("")

	// unpacked, err := Unpack("abccd")
	// fmt.Println("err", err)
	// fmt.Println("unpacked", unpacked)
	// fmt.Println("")

	// unpacked, err := Unpack("")
	// fmt.Println("err", err)
	// fmt.Println("unpacked", unpacked)
	// fmt.Println("")

	// unpacked, err := Unpack("aaa0b")
	// fmt.Println("err", err)
	// fmt.Println("unpacked", unpacked)
	// fmt.Println("")
}
