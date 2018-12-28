package ast

import (
	"dmeijboom/config/tokens"
)

type Node interface {
	Loc() *tokens.Location
}

type Stmt interface {
	Node
	stmtNode()
}
