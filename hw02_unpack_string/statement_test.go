package hw02unpackstring

import (
	"fmt"
	"testing"
)

// Test by a1b2c9de0 string.
var a1Lexeme, b2Lexeme, c9Lexeme, dLexeme, e0Lexeme Lexeme
var msg string

func TestMain(m *testing.M) {

	fmt.Println("Iterate test for \"a1b2c9de0\" unpacked.")

	fmt.Println("Lexemes for Statements blocks:")

	a1Lexeme.SetRune('a')
	a1Lexeme.SetCount(1)
	msg = fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", a1Lexeme.GetRune(), a1Lexeme.GetCount(), a1Lexeme.Unpack())
	fmt.Println(msg)

	b2Lexeme.SetRune('b')
	b2Lexeme.SetCount(2)
	msg = fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", b2Lexeme.GetRune(), b2Lexeme.GetCount(), b2Lexeme.Unpack())
	fmt.Println(msg)

	c9Lexeme.SetRune('c')
	c9Lexeme.SetCount(9)
	msg = fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", c9Lexeme.GetRune(), c9Lexeme.GetCount(), c9Lexeme.Unpack())
	fmt.Println(msg)

	dLexeme.SetRune('d')
	msg = fmt.Sprintf("Lexeme %q without any count unpacked to: %q.", dLexeme.GetRune(), dLexeme.Unpack())
	fmt.Println(msg)

	e0Lexeme.SetRune('e')
	e0Lexeme.SetCount(0)
	msg = fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", e0Lexeme.GetRune(), e0Lexeme.GetCount(), e0Lexeme.Unpack())
	fmt.Println(msg)

	m.Run()
}

func TestStatementUnpack(t *testing.T) {

	fmt.Println("Reverse realization.")

	// This is .......e0 block.

	var nilSatetment *StatementBlock

	var statementBlock_e0 = StatementBlock{
		BlockLexeme: e0Lexeme,
		NextBlock:   nilSatetment,
	}

	e0 := statementBlock_e0.Unpack()
	msg = fmt.Sprintf("Statement \"e0\" unpacked to: %q.", e0)
	fmt.Println(msg)
	expected_e0 := ""

	if e0 != "" {
		t.Errorf("(Statement \"e0\").Unpack() = %q, but expected %q.", e0, expected_e0)
	}

	// This is ......de0 blocks.

	var statementBlock_de0 = StatementBlock{
		BlockLexeme: dLexeme,
		NextBlock:   &statementBlock_e0,
	}
	de0 := statementBlock_de0.Unpack()
	msg = fmt.Sprintf("Statement \"de0\" unpacked to: %q.", de0)
	fmt.Println(msg)
	expected_de0 := "d"

	if de0 != expected_de0 {
		t.Errorf("(Statement \"de0\").Unpack() = %q, but expected %q.", de0, expected_de0)
	}

	// This is ....c9de0 blocks.

	var statementBlock_c9de0 = StatementBlock{
		BlockLexeme: c9Lexeme,
		NextBlock:   &statementBlock_de0,
	}
	c9de0 := statementBlock_c9de0.Unpack()
	msg = fmt.Sprintf("Statement \"c9de0\" unpacked to: %q.", c9de0)
	fmt.Println(msg)
	expected_c9de0 := "cccccccccd"

	if c9de0 != expected_c9de0 {
		t.Errorf("(Statement \"c9de0\").Unpack() = %q, but expected %q.", c9de0, expected_c9de0)
	}

	// This is ..b2c9de0 blocks.

	var statementBlock_b2c9de0 = StatementBlock{
		BlockLexeme: b2Lexeme,
		NextBlock:   &statementBlock_c9de0,
	}
	b2c9de0 := statementBlock_b2c9de0.Unpack()
	msg = fmt.Sprintf("Statement \"b2c9de0\" unpacked to: %q.", b2c9de0)
	fmt.Println(msg)
	expected_b2c9de0 := "bbcccccccccd"

	if b2c9de0 != expected_b2c9de0 {
		t.Errorf("(Statement \"b2c9de0\").Unpack() = %q, but expected %q.", b2c9de0, expected_b2c9de0)
	}

	// This is a1b2c9de0 blocks.

	var statementBlock_a1b2c9de0 = StatementBlock{
		BlockLexeme: a1Lexeme,
		NextBlock:   &statementBlock_b2c9de0,
	}

	a1b2c9de0 := statementBlock_a1b2c9de0.Unpack()
	msg = fmt.Sprintf("Statement \"a1b2c9de0\" unpacked to: %q.", a1b2c9de0)
	fmt.Println(msg)
	expected_a1b2c9de0 := "abbcccccccccd"
	if a1b2c9de0 != expected_a1b2c9de0 {
		t.Errorf("(Statement \"a1b2c9de0\").Unpack() = %q, but expected %q.", a1b2c9de0, expected_a1b2c9de0)
	}
}
