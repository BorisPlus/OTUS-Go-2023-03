package hw02unpackstring

type StatementBlock struct {
	BlockLexeme Lexeme
	NextBlock   *StatementBlock
}

func (statementBlock *StatementBlock) GetNext() *StatementBlock {
	return statementBlock.NextBlock
}

func (statementBlock *StatementBlock) GetLexeme() Lexeme {
	return statementBlock.BlockLexeme
}

func (statementBlock *StatementBlock) Unpack() string {
	if statementBlock == nil {
		return ""
	}
	nextStatementBlock := statementBlock.GetNext()
	lexeme := statementBlock.GetLexeme()
	return lexeme.Unpack() + nextStatementBlock.Unpack()
}
