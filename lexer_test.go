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
	assert.Equal(t, len(actual), len(expected), "Tokens length doesn't match")

	for i := 0; i < len(actual); i++ {
		tok_actual := actual[i]
		tok_expected := expected[i]

		assert.Equal(t, tok_actual.Kind, tok_expected.Kind, "Token kind doesn't match")
		assert.Equal(t, tok_actual.Value, tok_expected.Value, "Token value doesn't match")
	}
}

func TestLiterals(t *testing.T) {
	lexCmp(t, "\"This is a string\"", []tokens.Token{{tokens.String, "This is a string", nil}})
	lexCmp(t, "\"This is a\n \\\"string\"", []tokens.Token{{tokens.String, "This is a\n \"string", nil}})
	lexCmp(t, "true", []tokens.Token{{tokens.Boolean, true, nil}})
	lexCmp(t, "false", []tokens.Token{{tokens.Boolean, false, nil}})
	lexCmp(t, "1929", []tokens.Token{{tokens.Integer, 1929, nil}})
	lexCmp(t, "124.29403", []tokens.Token{{tokens.Float, 124.29403, nil}})
}

func TestUnfinishedString(t *testing.T) {
	lexShouldErr(t, "\"Start of string", "Unfinished string should error")
}

func TestWhitespace(t *testing.T) {
	lexCmp(t, "29.1 	\"Hi\" \n4", []tokens.Token{
		{tokens.Float, 29.1, nil},
		{tokens.String, "Hi", nil},
		{tokens.Integer, 4, nil},
	})
}
