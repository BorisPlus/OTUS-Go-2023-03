package hw02unpackstring

import (
	"fmt"
	"testing"
)

// a1b2c9de0
var lexeme_a1, lexeme_b2, lexeme_c9, lexeme_d, lexeme_e0 Lexeme
var msg string

func TestMain(m *testing.M) {

	fmt.Println("Iterate test for \"a1b2c9de0\" unpacked.")

	fmt.Println("Lexemes for Statements blocks:")

	lexeme_a1.SetRune('a')
	lexeme_a1.SetCount(1)
	msg = fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", lexeme_a1.GetRune(), lexeme_a1.GetCount(), lexeme_a1.Unpack())
	fmt.Println(msg)

	lexeme_b2.SetRune('b')
	lexeme_b2.SetCount(2)
	msg = fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", lexeme_b2.GetRune(), lexeme_b2.GetCount(), lexeme_b2.Unpack())
	fmt.Println(msg)

	lexeme_c9.SetRune('c')
	lexeme_c9.SetCount(9)
	msg = fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", lexeme_c9.GetRune(), lexeme_c9.GetCount(), lexeme_c9.Unpack())
	fmt.Println(msg)

	lexeme_d.SetRune('d')
	msg = fmt.Sprintf("Lexeme %q without any count unpacked to: %q.", lexeme_d.GetRune(), lexeme_d.Unpack())
	fmt.Println(msg)

	lexeme_e0.SetRune('e')
	lexeme_e0.SetCount(0)
	msg = fmt.Sprintf("Lexeme %q*%d unpacked to: %q.", lexeme_e0.GetRune(), lexeme_e0.GetCount(), lexeme_e0.Unpack())
	fmt.Println(msg)

	m.Run()
}

func TestStatementUnpack(t *testing.T) {

	fmt.Println("Reverse realization.")

	// .......e0

	var nilSatetment *StatementBlock

	var statementBlock_e0 = StatementBlock{
		BlockLexeme: lexeme_e0,
		NextBlock:   nilSatetment,
	}

	e0 := statementBlock_e0.Unpack()
	msg = fmt.Sprintf("Statement \"e0\" unpacked to: %q.", e0)
	fmt.Println(msg)
	expected_e0 := ""

	if e0 != "" {
		t.Errorf("(Statement \"e0\").Unpack() = %q, but expected %q.", e0, expected_e0)
	}

	// ......de0

	var statementBlock_de0 = StatementBlock{
		BlockLexeme: lexeme_d,
		NextBlock:   &statementBlock_e0,
	}
	de0 := statementBlock_de0.Unpack()
	msg = fmt.Sprintf("Statement \"de0\" unpacked to: %q.", de0)
	fmt.Println(msg)
	expected_de0 := "d"

	if de0 != expected_de0 {
		t.Errorf("(Statement \"de0\").Unpack() = %q, but expected %q.", de0, expected_de0)
	}

	// ....c9de0

	var statementBlock_c9de0 = StatementBlock{
		BlockLexeme: lexeme_c9,
		NextBlock:   &statementBlock_de0,
	}
	c9de0 := statementBlock_c9de0.Unpack()
	msg = fmt.Sprintf("Statement \"c9de0\" unpacked to: %q.", c9de0)
	fmt.Println(msg)
	expected_c9de0 := "cccccccccd"

	if c9de0 != expected_c9de0 {
		t.Errorf("(Statement \"c9de0\").Unpack() = %q, but expected %q.", c9de0, expected_c9de0)
	}

	// ..b2c9de0

	var statementBlock_b2c9de0 = StatementBlock{
		BlockLexeme: lexeme_b2,
		NextBlock:   &statementBlock_c9de0,
	}
	b2c9de0 := statementBlock_b2c9de0.Unpack()
	msg = fmt.Sprintf("Statement \"b2c9de0\" unpacked to: %q.", b2c9de0)
	fmt.Println(msg)
	expected_b2c9de0 := "bbcccccccccd"

	if b2c9de0 != expected_b2c9de0 {
		t.Errorf("(Statement \"b2c9de0\").Unpack() = %q, but expected %q.", b2c9de0, expected_b2c9de0)
	}

	// a1b2c9de0

	var statementBlock_a1b2c9de0 = StatementBlock{
		BlockLexeme: lexeme_a1,
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
