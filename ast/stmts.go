package ast

import (
	"dmeijboom/config/tokens"
)

type Section struct {
	Name *Ident
	Body []Node
}

func (section *Section) Loc() *tokens.Location {
	return section.Name.Loc()
}


type Typedef struct {
	Name *Ident
	Type *Type
}

func (typedef *Typedef) Loc() *tokens.Location {
	return typedef.Name.Loc()
}


type Let struct {
	Name *Ident
	Type *Type
	Value Node
}

func (let *Let) Loc() *tokens.Location {
	return let.Name.Loc()
}

/**
 * Statement definitions
 */
func (section *Section) stmtNode() {}
func (typedef *Typedef) stmtNode() {}
func (let *Let) stmtNode() {}
