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


type Call struct {
	Args []Expr
	Callee Expr
	Location *tokens.Location
}

func (call *Call) Loc() *tokens.Location {
	return call.Location
}

func (call *Call) Accept(visitor Visitor) {
	if len(call.Args) > 0 {
		for _, arg := range call.Args {
			arg.Accept(visitor)
			visitor.VisitInlineExpr(arg)
		}
	}

	call.Callee.Accept(visitor)
	visitor.VisitInlineExpr(call.Callee)
	visitor.VisitCall(call)
}


type Member struct {
	Object Expr
	Field Expr
	Location *tokens.Location
}

func (member *Member) Loc() *tokens.Location {
	return member.Location
}

func (member *Member) Accept(visitor Visitor) {
	member.Object.Accept(visitor)
	visitor.VisitInlineExpr(member.Object)
	member.Field.Accept(visitor)
	visitor.VisitMember(member)
}


/**
 * Expression definitions
 */
func (ident *Ident) exprNode() {}
func (init *Initialize) exprNode() {}
func (literal *Literal) exprNode() {}
func (call *Call) exprNode() {}
func (member *Member) exprNode() {}
