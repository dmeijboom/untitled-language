package main

import (
	"fmt"
	"dmeijboom/config/ast"
	"dmeijboom/config/tokens"
)

type Parser struct {
	tokens []tokens.Token
	scope *Scope
	index int
}

func NewParser(tokens []tokens.Token) *Parser {
	return &Parser{tokens: tokens, scope: NewScope(nil)}
}

func (parser *Parser) openScope() {
	parser.scope = NewScope(parser.scope)
}

func (parser *Parser) closeScope() []ast.Node {
	parent := parser.scope.parent
	nodes := parser.scope.Body()
	parser.scope = parent
	return nodes
}

func (parser *Parser) hasTokens() bool {
	return parser.index <= len(parser.tokens)-1
}

func (parser *Parser) pushBack() {
	parser.index--
}

func (parser *Parser) accept(kind tokens.TokenKind) bool {
	if parser.tokens[parser.index].Kind == kind {
		parser.index++
		return true
	}

	return false
}

func (parser *Parser) expect(kind tokens.TokenKind, value ...interface{}) *tokens.Token {
	token := parser.tok()

	if !parser.accept(kind) {
		panic(fmt.Errorf("SyntaxError: expecting type `%s` not `%s`", kind, token.Kind))
	} else if len(value) > 0 && token.Value != value[0] {
		panic(fmt.Errorf("SyntaxError: expecting token value %#v not %#v", value[0], token.Value))
	}

	return token
}

func (parser *Parser) tok() *tokens.Token {
	return &parser.tokens[parser.index]
}

func (parser *Parser) next() *tokens.Token {
	token := parser.tokens[parser.index]
	parser.index++
	return &token
}

func (parser *Parser) wrapError(str string) error {
	return fmt.Errorf("%s at %d:%d", str, parser.tok().Loc.Line, parser.tok().Loc.Column)
}

func (parser *Parser) ident() *ast.Ident {
	tok := parser.expect(tokens.Ident)

	return &ast.Ident{
		Location: tok.Loc,
		Value: tok.Value.(string),
	}
}

func (parser *Parser) section() {
	ident := parser.ident()
	parser.expect(tokens.LBracket)
	parser.openScope()
	parser.parseGlobal()
	body := parser.closeScope()
	parser.expect(tokens.RBracket)
	
	parser.scope.Add(&ast.Section{
		Name: ident,
		Body: body,
	})
}

func (parser *Parser) parseType() *ast.Type {
	array := false

	if parser.accept(tokens.LSqrBracket) {
		parser.expect(tokens.RSqrBracket)
	}

	name := parser.ident()
	var fields []ast.Field

	if name.Value == "object" &&
		parser.accept(tokens.LBracket) {
		fields = []ast.Field{}

		for parser.accept(tokens.Ident) {
			parser.pushBack()
			fieldName := parser.ident()
			parser.expect(tokens.Colon)
			fieldType := parser.parseType()

			fields = append(fields, ast.Field{
				Name: fieldName,
				Type: fieldType,
			})
		}

		parser.expect(tokens.RBracket)
	}

	return &ast.Type{
		Name: name,
		Array: array,
		Fields: fields,
		Optional: parser.accept(tokens.Query),
	}
}

func (parser *Parser) let() {
	parser.expect(tokens.Keyword, "let")
	name := parser.ident()
	parser.expect(tokens.Colon)
	type_ := parser.parseType()

	parser.scope.Add(&ast.Let{
		Name: name,
		Type: type_,
		Value: nil,
	})
}

func (parser *Parser) typedef() {
	parser.expect(tokens.Keyword, "type")
	name := parser.ident()
	parser.expect(tokens.Colon)
	type_ := parser.parseType()

	parser.scope.Add(&ast.Typedef{
		Name: name,
		Type: type_,
	})
}

func (parser *Parser) stmt() {
	if parser.accept(tokens.Ident) {
		if parser.accept(tokens.LBracket) {
			parser.pushBack()
			parser.pushBack()
			parser.section()
			return
		}

		parser.pushBack()
	} else if parser.accept(tokens.Keyword) {
		parser.pushBack()
		token := parser.tok()
		matched := true

		switch token.Value.(string) {
		case "type":
			parser.typedef()
			break
		case "let":
			parser.let()
			break
		default:
			matched = false
			break
		}

		if matched {
			return
		}
	}

	panic(fmt.Errorf("SyntaxError: unexpected %s", parser.tok()))
}

func (parser *Parser) parseGlobal() {
	for parser.hasTokens() {
		if parser.scope.parent != nil &&
			parser.accept(tokens.RBracket) {
			parser.pushBack()
			return
		}

		parser.stmt()
	}
}

func (parser *Parser) Parse() (nodes []ast.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(error); ok {
				err = parser.wrapError(r.(error).Error())
			} else {
				panic(r)
			}
		}
	}()

	parser.parseGlobal()
	nodes = parser.closeScope()

	return
}
