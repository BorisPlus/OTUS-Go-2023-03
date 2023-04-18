package hw02unpackstring

import (
	"fmt"
	"testing"
)

// Test by a1b2c9de0 string.
var (
	a1Lexeme, b2Lexeme, c9Lexeme, dLexeme, e0Lexeme Lexeme
	msg                                             string
)

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

	e0StatementBlock := StatementBlock{
		BlockLexeme: e0Lexeme,
		NextBlock:   nilSatetment,
	}

	e0 := e0StatementBlock.Unpack()
	msg = fmt.Sprintf("Statement \"e0\" unpacked to: %q.", e0)
	fmt.Println(msg)
	e0Expected := ""

	if e0 != "" {
		t.Errorf("(Statement \"e0\").Unpack() = %q, but expected %q.", e0, e0Expected)
	}

	// This is ......de0 blocks.

	de0StatementBlock := StatementBlock{
		BlockLexeme: dLexeme,
		NextBlock:   &e0StatementBlock,
	}
	de0 := de0StatementBlock.Unpack()
	msg = fmt.Sprintf("Statement \"de0\" unpacked to: %q.", de0)
	fmt.Println(msg)
	de0Expected := "d"

	if de0 != de0Expected {
		t.Errorf("(Statement \"de0\").Unpack() = %q, but expected %q.", de0, de0Expected)
	}

	// This is ....c9de0 blocks.

	c9de0StatementBlock := StatementBlock{
		BlockLexeme: c9Lexeme,
		NextBlock:   &de0StatementBlock,
	}
	c9de0 := c9de0StatementBlock.Unpack()
	msg = fmt.Sprintf("Statement \"c9de0\" unpacked to: %q.", c9de0)
	fmt.Println(msg)
	c9de0Expected := "cccccccccd"

	if c9de0 != c9de0Expected {
		t.Errorf("(Statement \"c9de0\").Unpack() = %q, but expected %q.", c9de0, c9de0Expected)
	}

	// This is ..b2c9de0 blocks.

	b2c9de0StatementBlock := StatementBlock{
		BlockLexeme: b2Lexeme,
		NextBlock:   &c9de0StatementBlock,
	}
	b2c9de0 := b2c9de0StatementBlock.Unpack()
	msg = fmt.Sprintf("Statement \"b2c9de0\" unpacked to: %q.", b2c9de0)
	fmt.Println(msg)
	b2c9de0Expected := "bbcccccccccd"

	if b2c9de0 != b2c9de0Expected {
		t.Errorf("(Statement \"b2c9de0\").Unpack() = %q, but expected %q.", b2c9de0, b2c9de0Expected)
	}

	// This is a1b2c9de0 blocks.

	a1b2c9de0StatementBlock := StatementBlock{
		BlockLexeme: a1Lexeme,
		NextBlock:   &b2c9de0StatementBlock,
	}

	a1b2c9de0 := a1b2c9de0StatementBlock.Unpack()
	msg = fmt.Sprintf("Statement \"a1b2c9de0\" unpacked to: %q.", a1b2c9de0)
	fmt.Println(msg)
	a1b2c9de0Expected := "abbcccccccccd"
	if a1b2c9de0 != a1b2c9de0Expected {
		t.Errorf("(Statement \"a1b2c9de0\").Unpack() = %q, but expected %q.", a1b2c9de0, a1b2c9de0Expected)
	}
}
