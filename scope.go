package main

import "dmeijboom/config/ast"

type Scope struct {
	parent *Scope
	body []ast.Node
}

func NewScope(parent *Scope) *Scope {
	return &Scope{parent: parent, body: []ast.Node{}}
}

func (scope *Scope) Add(node ast.Node) {
	scope.body = append(scope.body, node)
}

func (scope *Scope) Body() []ast.Node {
	return scope.body
}
