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

func (block *Block) Accept(visitor Visitor) {
	if len(block.Body) > 0 {
		for _, node := range block.Body {
			node.Accept(visitor)
		}
	}

	visitor.VisitBlock(block)
}


type Section struct {
	Name *Ident
	Block *Block
}

func (section *Section) Loc() *tokens.Location {
	return section.Name.Loc()
}

func (section *Section) Accept(visitor Visitor) {
	visitor.VisitSection(section)
	section.Name.Accept(visitor)
	section.Block.Accept(visitor)
}


type Typedef struct {
	Name *Ident
	Type *Type
}

func (typedef *Typedef) Loc() *tokens.Location {
	return typedef.Name.Loc()
}

func (typedef *Typedef) Accept(visitor Visitor) {
	typedef.Type.Accept(visitor)
	typedef.Name.Accept(visitor)
	visitor.VisitTypedef(typedef)
}


type Assign struct {
	Name *Ident
	Type *Type
	Value Node
}

func (assign *Assign) Loc() *tokens.Location {
	return assign.Name.Loc()
}

func (assign *Assign) Accept(visitor Visitor) {
	assign.Type.Accept(visitor)
	
	if assign.Value != nil {
		assign.Value.Accept(visitor)
	}

	visitor.VisitAssign(assign)
}

/**
 * Statement definitions
 */
func (section *Section) stmtNode() {}
func (typedef *Typedef) stmtNode() {}
func (block *Block) stmtNode() {}
func (assign *Assign) stmtNode() {}
