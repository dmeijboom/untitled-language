package main

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
	"dmeijboom/config/tokens"
)

var keywords = []string{
	"type", "let",
}

type Lexer struct {
	pos int
	line int
	col int
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

	if r == '\n' {
		lexer.line++
		lexer.col = 0
	} else {
		lexer.col += size
	}

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

func (lexer *Lexer) number() (tokens.Token, error) {
	num := ""
	is_float := false

	for !lexer.eof() {
		current := lexer.current()

		if !is_float && current == '.' {
			is_float = true
		} else if !unicode.IsNumber(current) {
			break
		}

		num += string(lexer.next())
	}
	
	if is_float {
		floatval, err := strconv.ParseFloat(num, 64)
		return tokens.Token{
			Kind: tokens.Float,
			Value: floatval,
		}, err
	}

	numval, err := strconv.Atoi(num)
	return tokens.Token{
		Kind: tokens.Integer,
		Value: numval,
	}, err
}

func (lexer *Lexer) string() (tokens.Token, error) {
	str := ""

	lexer.next()

	loop:
	for {
		if lexer.eof() {
			return tokens.Token{}, fmt.Errorf("Unfinished string literal at %d:%d", lexer.line, lexer.pos)
		}

		current := lexer.current()

		switch current {
		case '"':
			break loop
		case '\\':
			lexer.next()
			str += string(lexer.next())
			continue loop
		default:
			str += string(current)
			lexer.next()
		}
	}

	lexer.next()

	return tokens.Token{
		Kind: tokens.String,
		Value: str,
	}, nil
}

func (lexer *Lexer) isKeyword(ident string) bool {
	for _, keyword := range keywords {
		if keyword == ident {
			return true
		}
	}

	return false
}

func (lexer *Lexer) Lex() ([]tokens.Token, error) {
	tokenList := []tokens.Token{}

	loop:
	for !lexer.eof() {
		current := lexer.current()
		start_pos := lexer.col

		var token tokens.Token
		var err error

		switch current {
		case '"':
			token, err = lexer.string()
			break
		case '=':
			token = tokens.Token{Kind: tokens.Equals}
			lexer.next()
			break
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
			continue loop
		case ' ', '\t', '\r':
			lexer.next()
			continue loop
		default:
			if unicode.IsNumber(current) {
				token, err = lexer.number()
				break
			} else if unicode.IsLetter(current) {
				token, err = lexer.ident()

				if err == nil &&
					(token.Value == "true" || token.Value == "false") {
					token.Kind = tokens.Boolean
					token.Value = token.Value == "true"
				} else if lexer.isKeyword(token.Value.(string)) {
					token.Kind = tokens.Keyword
				}
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
