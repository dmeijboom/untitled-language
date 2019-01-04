package compiler

import (
	"dmeijboom/config/ast"
)

var builtinTypes = []string{
	"int", "bool", "string", "float",
}

type Compiler struct {
	source *ast.Source
	instructions []Instruction
}

func NewCompiler(source *ast.Source) *Compiler {
	return &Compiler{
		source: source,
		instructions: []Instruction{},
	}
}

func (compiler *Compiler) add(instruction Instruction) {
	compiler.instructions = append(compiler.instructions, instruction)
}

func (compiler *Compiler) isBuiltin(name string) bool {
	for _, typeName := range builtinTypes {
		if typeName == name {
			return true
		}
	}

	return false
}

func (compiler *Compiler) VisitIdent(ident *ast.Ident) {

}

func (compiler *Compiler) VisitField(field *ast.Field) {
	compiler.add(&MakeField{
		Name: field.Name.Value,
		Location: field.Loc(),
	})
}

func (compiler *Compiler) VisitType(node *ast.Type) {
	if node.Name.Value == "object" {
		compiler.add(&MakeObject{
			Fields: len(node.Fields),
			Location: node.Loc(),
		})
		return
	}

	loadType := &LoadType{
		Array: node.Array,
		Optional: node.Optional,
		Location: node.Loc(),
	}

	if compiler.isBuiltin(node.Name.Value) {
		switch node.Name.Value {
		case "int":
			loadType.Type = IntegerType
			break
		case "bool":
			loadType.Type = BooleanType
			break
		case "string":
			loadType.Type = StringType
			break
		case "float":
			loadType.Type = FloatType
			break
		}
	} else {
		loadType.Type = UserType
		loadType.TypeName = node.Name.Value
	}

	compiler.add(loadType)
}

func (compiler *Compiler) VisitPreInitialize(init *ast.Initialize) {
	compiler.add(&NewObject{
		Fields: len(init.Fields),
		Location: init.Loc(),
	})
}

func (compiler *Compiler) VisitInitialize(init *ast.Initialize) {
	compiler.add(&Initialize{
		Location: init.Loc(),
	})
}

func (compiler *Compiler) VisitInitializeField(initField *ast.InitializeField) {
	compiler.add(&SetField{
		Name: initField.Name.Value,
		Location: initField.Loc(),
	})
}

func (compiler *Compiler) VisitLiteral(literal *ast.Literal) {
	compiler.add(&LoadConst{
		Value: literal.Value,
		Location: literal.Loc(),
	})
}

func (compiler *Compiler) VisitBlock(block *ast.Block) {

}

func (compiler *Compiler) VisitSection(section *ast.Section) {
	compiler.add(&SetSection{
		Name: section.Name.Value,
		Location: section.Loc(),
	})
}

func (compiler *Compiler) VisitTypedef(typedef *ast.Typedef) {
	compiler.add(&MakeType{
		Name: typedef.Name.Value,
		Location: typedef.Loc(),
	})
}

func (compiler *Compiler) VisitAssign(assign *ast.Assign) {
	compiler.add(&StoreVal{
		Name: assign.Name.Value,
		HasValue: assign.Value != nil,
		Location: assign.Loc(),
	})
}

func (compiler *Compiler) VisitSource(source *ast.Source) {

}

func (compiler *Compiler) Compile() ([]Instruction, error) {
	compiler.source.Accept(compiler)
	return compiler.instructions, nil
}
