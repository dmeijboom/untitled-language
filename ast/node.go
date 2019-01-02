package ast

import (
	"dmeijboom/config/tokens"
)

type Node interface {
	Visitable
	Loc() *tokens.Location
}

type Stmt interface {
	Node
	stmtNode()
}

type Source struct {
	Filename string
	Block *Block
}

func (source *Source) Accept(visitor Visitor) {
	source.Block.Accept(visitor)
	visitor.VisitSource(source)
}
