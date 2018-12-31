package main

import (
	"reflect"
	"testing"
	"dmeijboom/config/ast"
	"github.com/stretchr/testify/assert"
)

func parseCmpNode(t *testing.T, a ast.Node, b ast.Node) {
	if !assert.Equal(t, a == nil, b == nil, "AST nodes should be both nil or both not nil") ||
		(a == nil && b == nil) {
		return
	}

	type_a := reflect.TypeOf(a).Elem().Name()
	type_b := reflect.TypeOf(b).Elem().Name()

	assert.Equal(t, type_a, type_b, "AST nodes differ")

	switch node_a := a.(type) {
	case *ast.Block:
		node_b := b.(*ast.Block)
		assert.Equal(t, len(node_a.Body), len(node_b.Body), "Block body count doesn't match")

		for i := 0; i < len(node_a.Body); i++ {
			parseCmpNode(t, node_a.Body[i], node_b.Body[i])
		}
		break
	case *ast.Section:
		node_b := b.(*ast.Section)
		parseCmpNode(t, node_a.Name, node_b.Name)
		parseCmpNode(t, node_a.Block, node_b.Block)
		break
	case *ast.Typedef:
		node_b := b.(*ast.Typedef)
		parseCmpNode(t, node_a.Name, node_b.Name)
		parseCmpNode(t, node_a.Type, node_b.Type)
		break
	case *ast.Literal:
		node_b := b.(*ast.Literal)
		assert.Equal(t, node_a.Type, node_b.Type, "Literal type doesn't match")
		assert.Equal(t, node_a.Value, node_b.Value, "Literal value doesn't match")
		break
	case *ast.Ident:
		node_b := b.(*ast.Ident)
		assert.Equal(t, node_a.Value, node_b.Value, "Ident value doesn't match")
		break
	case *ast.Type:
		node_b := b.(*ast.Type)
		parseCmpNode(t, node_a.Name, node_b.Name)
		assert.Equal(t, node_a.Array, node_b.Array, "Type should be an array")
		assert.Equal(t, node_a.Optional, node_b.Optional, "Type should be optional")

		assert.Equal(t, len(node_a.Fields), len(node_b.Fields), "Type field-length doesn't match")

		for i := 0; i < len(node_a.Fields); i++ {
			parseCmpNode(t, &node_a.Fields[i], &node_b.Fields[i])
		}
		break
	case *ast.Field:
		node_b := b.(*ast.Field)
		parseCmpNode(t, node_a.Name, node_b.Name)
		parseCmpNode(t, node_a.Type, node_b.Type)
		break
	case *ast.Let:
		node_b := b.(*ast.Let)
		parseCmpNode(t, node_a.Name, node_b.Name)
		parseCmpNode(t, node_a.Type, node_b.Type)
		parseCmpNode(t, node_a.Value, node_b.Value)
		break
	}
}

func parseCmp(t *testing.T, input string, expected []ast.Node) {
	lexer := NewLexer(input)
	tokens, err := lexer.Lex()

	assert.Nil(t, err)

	parser := NewParser(tokens)
	source, err := parser.Parse()

	assert.Nil(t, err)
	assert.Equal(t, len(source.Block.Body), len(expected), "AST Node length doesn't match")

	for i := 0; i < len(source.Block.Body); i++ {
		parseCmpNode(t, source.Block.Body[i], expected[i])
	}
}

func TestSection(t *testing.T) {
	parseCmp(t, "testSection {}", []ast.Node{&ast.Section{
		Name: &ast.Ident{"testSection", nil},
		Block: &ast.Block{[]ast.Node{}, nil},
	}})

	parseCmp(t, `testSection {

	} example {}`, []ast.Node{
		&ast.Section{&ast.Ident{"testSection", nil}, &ast.Block{[]ast.Node{}, nil}},
		&ast.Section{&ast.Ident{"example", nil}, &ast.Block{[]ast.Node{}, nil}},
	})
}

func TestTypedef(t *testing.T) {
	parseCmp(t, `testSection {
		type Test: int
	}`, []ast.Node{&ast.Section{
		Name: &ast.Ident{"testSection", nil},
		Block: &ast.Block{[]ast.Node{
			&ast.Typedef{
				Name: &ast.Ident{"Test", nil},
				Type: &ast.Type{
					Array: false,
					Optional: false,
					Name: &ast.Ident{"int", nil},
				},
			},
		}, nil},
	}})

	parseCmp(t, `testSection {
		type User: object {
			name: string
			age: int
			email: string?
		}
	}`, []ast.Node{&ast.Section{
		&ast.Ident{"testSection", nil},
		&ast.Block{[]ast.Node{
			&ast.Typedef{
				Name: &ast.Ident{"User", nil},
				Type: &ast.Type{
					Name: &ast.Ident{"object", nil},
					Fields: []ast.Field{
						ast.Field{
							Name: &ast.Ident{"name", nil},
							Type: &ast.Type{Name: &ast.Ident{"string", nil}},
						},
						ast.Field{
							Name: &ast.Ident{"age", nil},
							Type: &ast.Type{Name: &ast.Ident{"int", nil}},
						},
						ast.Field{
							Name: &ast.Ident{"email", nil},
							Type: &ast.Type{Name: &ast.Ident{"string", nil}, Optional: true},
						},
					},
				},
			},
		}, nil},
	}})
}

func TestLet(t *testing.T) {
	parseCmp(t, `testSection { let name: string }`, []ast.Node{&ast.Section{
		&ast.Ident{"testSection", nil},
		&ast.Block{[]ast.Node{&ast.Let{
			Name: &ast.Ident{"name", nil},
			Type: &ast.Type{
				Name: &ast.Ident{"string", nil},
			},
		}}, nil},
	}})

	parseCmp(t, `testSection { let name: string = "Hello World" }`, []ast.Node{&ast.Section{
		&ast.Ident{"testSection", nil},
		&ast.Block{[]ast.Node{&ast.Let{
			Name: &ast.Ident{"name", nil},
			Type: &ast.Type{Name: &ast.Ident{"string", nil}},
			Value: &ast.Literal{ast.String, "Hello World", nil},
		}}, nil},
	}})
}
