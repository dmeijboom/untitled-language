package ast

import (
	"dmeijboom/config/tokens"
)

type Ident struct {
	Location *tokens.Location
	Value string
}

func (ident *Ident) Loc() *tokens.Location {
	return ident.Location
}


type Field struct {
	Name *Ident
	Type *Type
}

func (field *Field) Loc() *tokens.Location {
	return field.Name.Loc()
}


type Type struct {
	Name *Ident
	Array bool
	Optional bool
	Fields []Field
}

func (type_ *Type) Loc() *tokens.Location {
	return type_.Name.Loc()
}


type LiteralType int

const (
	String LiteralType = iota
	Integer
	Float
	Boolean
)

type Literal struct {
	Location *tokens.Location
	Type LiteralType
	Value interface{}
}

func (literal *Literal) Loc() *tokens.Location {
	return literal.Location
}
