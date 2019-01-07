package main

import (
	"testing"
	"dmeijboom/config/tokens"
	"github.com/stretchr/testify/assert"
)

func lexShouldErr(t *testing.T, input string, msgAndArgs ...interface{}) {
	lexer := NewLexer(input)
	_, err := lexer.Lex()

	assert.NotNil(t, err, msgAndArgs...)
}

func lexCmp(t *testing.T, input string, expected []tokens.Token) {
	lexer := NewLexer(input)
	actual, err := lexer.Lex()

	assert.Nil(t, err)
	assert.Equal(t, len(expected), len(actual), "Tokens length doesn't match")

	for i := 0; i < len(actual); i++ {
		tok_actual := actual[i]
		tok_expected := expected[i]

		assert.Equal(t, tok_expected.Kind, tok_actual.Kind, "Token kind doesn't match")
		assert.Equal(t, tok_expected.Value, tok_actual.Value, "Token value doesn't match")
	}
}

func TestLiterals(t *testing.T) {
	lexCmp(t, "\"This is a string\"", []tokens.Token{
		{tokens.String, "This is a string", nil},
		{tokens.EndStmt, nil, nil},
	})
	lexCmp(t, "\"This is a\n \\\"string\"", []tokens.Token{
		{tokens.String, "This is a\n \"string", nil},
		{tokens.EndStmt, nil, nil},
	})
	lexCmp(t, "true", []tokens.Token{
		{tokens.Boolean, true, nil},
		{tokens.EndStmt, nil, nil},
	})
	lexCmp(t, "false", []tokens.Token{
		{tokens.Boolean, false, nil},
		{tokens.EndStmt, nil, nil},
	})
	lexCmp(t, "1929", []tokens.Token{
		{tokens.Integer, 1929, nil},
		{tokens.EndStmt, nil, nil},
	})
	lexCmp(t, "124.29403", []tokens.Token{
		{tokens.Float, 124.29403, nil},
		{tokens.EndStmt, nil, nil},
	})
}

func TestUnfinishedString(t *testing.T) {
	lexShouldErr(t, "\"Start of string", "Unfinished string should error")
}

func TestWhitespace(t *testing.T) {
	lexCmp(t, "29.1 	\"Hi\" \n4", []tokens.Token{
		{tokens.Float, 29.1, nil},
		{tokens.String, "Hi", nil},
		{tokens.EndStmt, nil, nil},
		{tokens.Integer, 4, nil},
		{tokens.EndStmt, nil, nil},
	})
}

func TestIdent(t *testing.T) {
	lexCmp(t, "1name", []tokens.Token{
		{tokens.Integer, 1, nil},
		{tokens.Ident, "name", nil},
		{tokens.EndStmt, nil, nil},
	})

	lexCmp(t, "name1_testExample", []tokens.Token{
		{tokens.Ident, "name1_testExample", nil},
		{tokens.EndStmt, nil, nil},
	})

	lexShouldErr(t, "name-example", "Cannot use `-` in identifier")
}

func TestAutoEndStmt(t *testing.T) {
	lexCmp(t, `type name: string
	type age:
	int
	let x: Example {
	}`, []tokens.Token{
		{tokens.Keyword, "type", nil},
		{tokens.Ident, "name", nil},
		{tokens.Colon, nil, nil},
		{tokens.Ident, "string", nil},
		{tokens.EndStmt, nil, nil},

		{tokens.Keyword, "type", nil},
		{tokens.Ident, "age", nil},
		{tokens.Colon, nil, nil},
		{tokens.Ident, "int", nil},
		{tokens.EndStmt, nil, nil},

		{tokens.Keyword, "let", nil},
		{tokens.Ident, "x", nil},
		{tokens.Colon, nil, nil},
		{tokens.Ident, "Example", nil},
		{tokens.LBracket, nil, nil},
		{tokens.RBracket, nil, nil},
		{tokens.EndStmt, nil, nil},
	})
}
