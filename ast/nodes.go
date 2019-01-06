package ast

import (
	"dmeijboom/config/tokens"
)

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

	type_.Name.Accept(visitor)
	visitor.VisitType(type_)
}


type InitializeField struct {
	Name *Ident
	Value Expr
	Location *tokens.Location
}

func (initField *InitializeField) Loc() *tokens.Location {
	return initField.Location
}

func (initField *InitializeField) Accept(visitor Visitor) {
	initField.Name.Accept(visitor)
	initField.Value.Accept(visitor)
	visitor.VisitInlineExpr(initField.Value)
	visitor.VisitInitializeField(initField)
}
