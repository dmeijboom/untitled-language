package main

import (
	"fmt"
	"errors"
	"strings"
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

func (parser *Parser) accept(kind tokens.TokenKind, value ...interface{}) bool {
	token := parser.tok()

	if token != nil &&
		token.Kind == kind {
		if len(value) > 0 && token.Value != value[0] {
			return false
		}

		parser.index++
		return true
	}

	return false
}

func (parser *Parser) expect(kind tokens.TokenKind, value ...interface{}) *tokens.Token {
	token := parser.tok()

	if !parser.hasTokens() {
		panic(fmt.Errorf("SyntaxError: Unexpected EOF"))
	} else if !parser.accept(kind) {
		panic(fmt.Errorf("SyntaxError: expecting type `%s` not %s", kind, token.String()))
	} else if len(value) > 0 && token.Value != value[0] {
		panic(fmt.Errorf("SyntaxError: expecting token value %#v not %#v", value[0], token.Value))
	}

	return token
}

func (parser *Parser) tok() *tokens.Token {
	if !parser.hasTokens() {
		return nil
	}

	return &parser.tokens[parser.index]
}

func (parser *Parser) next() *tokens.Token {
	token := parser.tokens[parser.index]
	parser.index++
	return &token
}

func (parser *Parser) wrapError(str string) error {
	token := parser.tok()

	if !parser.hasTokens() {
		token = &parser.tokens[len(parser.tokens)-1]
	}

	return fmt.Errorf("%s at %d:%d", str, token.Loc.Line, token.Loc.Column)
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
	loc := parser.tok().Loc
	parser.parseGlobal()
	body := parser.closeScope()
	parser.expect(tokens.RBracket)
	parser.expect(tokens.EndStmt)
	
	parser.scope.Add(&ast.Section{
		Name: ident,
		Block: &ast.Block{Body: body, Location: loc},
	})
}

func (parser *Parser) objectdef() *ast.Type {
	parser.expect(tokens.Ident, "object")
	parser.pushBack()

	name := parser.ident()
	fields := []ast.Field{}

	parser.accept(tokens.LBracket)

	for parser.accept(tokens.Ident) {
		parser.pushBack()
		fieldName := parser.ident()
		parser.expect(tokens.Colon)
		fieldType := parser.parseType()
		parser.expect(tokens.EndStmt)

		fields = append(fields, ast.Field{
			Name: fieldName,
			Type: fieldType,
		})
	}

	parser.expect(tokens.RBracket)

	return &ast.Type{
		Name: name,
		Fields: fields,
	}
}

func (parser *Parser) parseType() *ast.Type {
	array := false

	if parser.accept(tokens.LSqrBracket) {
		parser.expect(tokens.RSqrBracket)
		array = true
	}

	name := parser.ident()

	if name.Value == "object" {
		panic(errors.New("SyntaxError: Cannot use object type outside typedef"))
	}

	return &ast.Type{
		Name: name,
		Array: array,
		Optional: parser.accept(tokens.Query),
	}
}

func (parser *Parser) init() *ast.Initialize {
	new := parser.expect(tokens.Keyword, "new")
	parser.expect(tokens.LBracket)

	fields := []ast.InitializeField{}

	for {
		if !parser.accept(tokens.Ident) {
			break
		}

		parser.pushBack()
		name := parser.ident()

		parser.expect(tokens.Equals)
		expr := parser.expr()
		parser.expect(tokens.EndStmt)

		fields = append(fields, ast.InitializeField{
			Name: name,
			Value: expr,
			Location: name.Loc(),
		})
	}

	parser.expect(tokens.RBracket)

	return &ast.Initialize{
		Fields: fields,
		Location: new.Loc,
	}
}

func (parser *Parser) member() ast.Expr {
	object := parser.ident()
	parser.expect(tokens.Interpunct)
	field := parser.ident()

	return &ast.Member{
		Object: object,
		Field: field,
		Location: object.Loc(),
	}
}

func (parser *Parser) call(callee ast.Expr) ast.Expr {
	parser.expect(tokens.LParent)

	args := []ast.Expr{}

	for !parser.accept(tokens.RParent) {
		args = append(args, parser.expr())
	}

	return &ast.Call{
		Args: args,
		Callee: callee,
		Location: callee.Loc(),
	}
}

func (parser *Parser) expr() ast.Expr {
	token := parser.tok()

	if parser.accept(tokens.Keyword, "new") {
		parser.pushBack()
		return parser.init()
	} else if parser.accept(tokens.Ident) {
		var expr ast.Expr

		if parser.accept(tokens.Interpunct) {
			parser.pushBack()
			parser.pushBack()
			expr = parser.member()
		} else {
			parser.pushBack()
			expr = parser.ident()
		}

		if parser.accept(tokens.LParent) {
			parser.pushBack()
			return parser.call(expr)
		}

		return expr
	} else if parser.accept(tokens.String) {
		return &ast.Literal{
			Type: ast.String,
			Value: token.Value,
		}
	} else if parser.accept(tokens.Integer) {
		return &ast.Literal{
			Type: ast.Integer,
			Value: token.Value,
		}
	} else if parser.accept(tokens.Float) {
		return &ast.Literal{
			Type: ast.Float,
			Value: token.Value,
		}
	} else if parser.accept(tokens.Boolean) {
		return &ast.Literal{
			Type: ast.Boolean,
			Value: token.Value,
		}
	}

	panic(fmt.Errorf("SyntaxError: unexpected %s", token))
}

func (parser *Parser) assign() {
	parser.expect(tokens.Keyword, "let")
	name := parser.ident()
	parser.expect(tokens.Colon)
	type_ := parser.parseType()

	var value ast.Expr

	if parser.accept(tokens.Equals) {
		value = parser.expr()
	}

	parser.expect(tokens.EndStmt)

	parser.scope.Add(&ast.Assign{
		Name: name,
		Type: type_,
		Value: value,
	})
}

func (parser *Parser) typedef() {
	parser.expect(tokens.Keyword, "type")
	name := parser.ident()
	parser.expect(tokens.Colon)

	var typeval *ast.Type

	if parser.accept(tokens.Ident, "object") {
		parser.pushBack()
		typeval = parser.objectdef()
	} else {
		typeval = parser.parseType()
	}

	parser.expect(tokens.EndStmt)

	parser.scope.Add(&ast.Typedef{
		Name: name,
		Type: typeval,
	})
}

func (parser *Parser) exprStmt() {
	expr := parser.expr()

	parser.expect(tokens.EndStmt)
	parser.scope.Add(&ast.ExprStmt{
		Expr: expr,
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
			parser.assign()
			break
		default:
			matched = false
			break
		}

		if matched {
			return
		}
	}
	
	parser.exprStmt()
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

func (parser *Parser) Parse() (source *ast.Source, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(error); ok && strings.Contains(r.(error).Error(), "SyntaxError:") {
				err = parser.wrapError(r.(error).Error())
			} else {
				panic(r)
			}
		}
	}()

	parser.parseGlobal()
	source = &ast.Source{Block: &ast.Block{
		Body: parser.closeScope(),
		Location: &tokens.Location{Line: 0, Column: 0},
	}}

	return
}
