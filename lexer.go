package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
	"dmeijboom/config/tokens"
)

type Lexer struct {
	pos int
	line int
	input string
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, line: 1}
}

func (lexer *Lexer) eof() bool {
	return lexer.pos >= len(lexer.input)
}

func (lexer *Lexer) current() rune {
	r, _ := utf8.DecodeRuneInString(lexer.input[lexer.pos:])
	return r
}

func (lexer *Lexer) next() rune {
	r, size := utf8.DecodeRuneInString(lexer.input[lexer.pos:])
	lexer.pos += size
	return r
}

func (lexer *Lexer) ident() (tokens.Token, error) {
	ident := ""

	for !lexer.eof() &&
		unicode.IsLetter(lexer.current()) {
		ident += string(lexer.next())
	}

	return tokens.Token{
		Kind: tokens.Ident,
		Value: ident,
	}, nil
}

func (lexer *Lexer) Lex() ([]tokens.Token, error) {
	tokenList := []tokens.Token{}

	loop:
	for !lexer.eof() {
		current := lexer.current()
		start_pos := lexer.pos

		var token tokens.Token
		var err error

		switch current {
		case '{':
			token = tokens.Token{Kind: tokens.LBracket}
			lexer.next()
			break
		case '}':
			token = tokens.Token{Kind: tokens.RBracket}
			lexer.next()
			break
		case ':':
			token = tokens.Token{Kind: tokens.Colon}
			lexer.next()
			break
		case '?':
			token = tokens.Token{Kind: tokens.Query}
			lexer.next()
			break
		case '[':
			token = tokens.Token{Kind: tokens.LSqrBracket}
			lexer.next()
			break
		case ']':
			token = tokens.Token{Kind: tokens.RSqrBracket}
			lexer.next()
			break
		case '\n':
			lexer.next()
			lexer.line++
			continue loop
		case ' ':
			lexer.next()
			continue loop
		default:
			if unicode.IsLetter(current) {
				token, err = lexer.ident()
				break
			}
			
			return nil, fmt.Errorf("Unknown token %q at %d:%d", lexer.current(), lexer.line, lexer.pos)
		}

		if err != nil {
			return nil, err
		}

		token.Loc = &tokens.Location{
			Line: lexer.line,
			Column: start_pos,
		}
		tokenList = append(tokenList, token)
	}

	return tokenList, nil
}
