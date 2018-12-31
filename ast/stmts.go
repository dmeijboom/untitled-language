package ast

import (
	"dmeijboom/config/tokens"
)

type Block struct {
	Body []Node
	Location *tokens.Location
}

func (block *Block) Loc() *tokens.Location {
	return block.Location
}


type Section struct {
	Name *Ident
	Block *Block
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
func (block *Block) stmtNode() {}
func (let *Let) stmtNode() {}
