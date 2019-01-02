package ast

import (
	"dmeijboom/config/tokens"
)

type Ident struct {
	Value string
	Location *tokens.Location
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


type InitializeField struct {
	Name *Ident
	Value Node
	Location *tokens.Location
}

func (initField *InitializeField) Loc() *tokens.Location {
	return initField.Location
}


type Initialize struct {
	Fields []InitializeField
	Location *tokens.Location
}

func (init *Initialize) Loc() *tokens.Location {
	return init.Location
}


type LiteralType int

const (
	String LiteralType = iota
	Integer
	Float
	Boolean
)

type Literal struct {
	Type LiteralType
	Value interface{}
	Location *tokens.Location
}

func (literal *Literal) Loc() *tokens.Location {
	return literal.Location
}
