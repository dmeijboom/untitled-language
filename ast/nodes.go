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

func (ident *Ident) Accept(visitor Visitor) {
	visitor.VisitIdent(ident)
}


type Field struct {
	Name *Ident
	Type *Type
}

func (field *Field) Loc() *tokens.Location {
	return field.Name.Loc()
}

func (field *Field) Accept(visitor Visitor) {
	field.Type.Accept(visitor)
	field.Name.Accept(visitor)
	visitor.VisitField(field)
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

func (type_ *Type) Accept(visitor Visitor) {
	if len(type_.Fields) > 0 {
		for _, field := range type_.Fields {
			field.Accept(visitor)
		}
	}

	visitor.VisitType(type_)
}


type InitializeField struct {
	Name *Ident
	Value Node
	Location *tokens.Location
}

func (initField *InitializeField) Loc() *tokens.Location {
	return initField.Location
}

func (initField *InitializeField) Accept(visitor Visitor) {
	initField.Name.Accept(visitor)
	initField.Value.Accept(visitor)
	visitor.VisitInitializeField(initField)
}


type Initialize struct {
	Fields []InitializeField
	Location *tokens.Location
}

func (init *Initialize) Loc() *tokens.Location {
	return init.Location
}

func (init *Initialize) Accept(visitor Visitor) {
	visitor.VisitPreInitialize(init)

	if len(init.Fields) > 0 {
		for _, field := range init.Fields {
			field.Accept(visitor)
		}
	}

	visitor.VisitInitialize(init)
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

func (literal *Literal) Accept(visitor Visitor) {
	visitor.VisitLiteral(literal)
}
