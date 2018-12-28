package main

import (
	"dmeijboom/config/ast"
	"dmeijboom/config/tokens"
)

type Parser struct {
	tokens []tokens.Token
	index int
}

func (parser *Parser) Parse() []ast.Stmt {
	return nil
}
