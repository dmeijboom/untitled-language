package tokens

import (
	"fmt"
)

type TokenKind int

const (
	Ident TokenKind = iota
	Keyword
	LBracket
	RBracket
	LParent
	RParent
	LSqrBracket
	RSqrBracket
	Colon
	Query
	Equals
	Interpunct
	String
	Boolean
	Integer
	Float
)

type Token struct {
	Kind TokenKind
	Value interface{}
	Loc *Location
}

func (token Token) String() string {
	if token.Value == nil {
		return fmt.Sprintf("(Token) %s", token.Kind.String())
	}

	return fmt.Sprintf("(Token) %s = %v", token.Kind.String(), token.Value)
}

func (kind TokenKind) String() string {
	switch kind {
	case Ident:
		return "Ident"
	case LBracket:
		return "LBracket"
	case RBracket:
		return "RBracket"
	case LParent:
		return "LParent"
	case RParent:
		return "RParent"
	case LSqrBracket:
		return "LSqrBracket"
	case RSqrBracket:
		return "RSqrBracket"
	case Colon:
		return "Colon"
	case Query:
		return "Query"
	case Equals:
		return "Equals"
	case String:
		return "String"
	case Boolean:
		return "Boolean"
	case Integer:
		return "Integer"
	case Float:
		return "Float"
	case Keyword:
		return "Keyword"
	case Interpunct:
		return "Interpunct"
	default:
		panic("Unknown token kind")
	}
}
