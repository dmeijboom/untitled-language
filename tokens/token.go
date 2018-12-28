package tokens

import (
	"fmt"
)

type TokenKind int

const (
	Ident TokenKind = iota
	LBracket
	RBracket
	LParent
	RParent
	LSqrBracket
	RSqrBracket
	Colon
	Query
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
	var name string

	switch token.Kind {
	case Ident:
		name = "Ident"
	case LBracket:
		name = "LBracket"
	case RBracket:
		name = "RBracket"
	case LParent:
		name = "LParent"
	case RParent:
		name = "RParent"
	case LSqrBracket:
		name = "LSqrBracket"
	case RSqrBracket:
		name = "RSqrBracket"
	case Colon:
		name = "Colon"
	case Query:
		name = "Query"
	case String:
		name = "String"
	case Boolean:
		name = "Boolean"
	case Integer:
		name = "Integer"
	case Float:
		name = "Float"
	}

	if token.Value == nil {
		return fmt.Sprintf("(Token) %s", name)
	}

	return fmt.Sprintf("(Token) %s = %v", name, token.Value)
}
