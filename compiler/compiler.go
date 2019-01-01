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

func (compiler *Compiler) compile(nodeInterface ast.Node) error {
	switch node := nodeInterface.(type) {
	case *ast.Block:
		for _, subNode := range node.Body {
			if err := compiler.compile(subNode); err != nil {
				return err
			}
		}
		break
	case *ast.Field:
		if err := compiler.compile(node.Type); err != nil {
			return err
		}

		compiler.add(&MakeField{
			Name: node.Name.Value,
			Location: node.Loc(),
		})
		break
	case *ast.Section:
		compiler.add(&SetSection{
			Name: node.Name.Value,
			Location: node.Loc(),
		})

		if err := compiler.compile(node.Block); err != nil {
			return err
		}
		break
	case *ast.Type:
		makeType := &MakeType{
			Array: node.Array,
			Fields: len(node.Fields),
			Optional: node.Optional,
			Location: node.Loc(),
		}

		if len(node.Fields) > 0 {
			for _, field := range node.Fields {
				if err := compiler.compile(&field); err != nil {
					return err
				}
			}
		}

		if compiler.isBuiltin(node.Name.Value) {
			switch node.Name.Value {
			case "int":
				makeType.Type = IntegerType
				break
			case "bool":
				makeType.Type = BooleanType
				break
			case "string":
				makeType.Type = StringType
				break
			case "float":
				makeType.Type = FloatType
				break
			case "object":
				makeType.Type = ObjectType
				break
			}
		} else {
			makeType.Type = UserType
			makeType.TypeName = node.Name.Value
		}

		compiler.add(makeType)
		break
	case *ast.Typedef:
		if err := compiler.compile(node.Type); err != nil {
			return err
		}

		compiler.add(&StoreVal{
			Name: node.Name.Value,
			Location: node.Loc(),
		})
		break
	case *ast.Assign:
		if err := compiler.compile(node.Type); err != nil {
			return err
		}

		compiler.add(&StoreVal{
			Name: node.Name.Value,
			Location: node.Loc(),
		})
		break
	}

	return nil
}

func (compiler *Compiler) Compile() ([]Instruction, error) {
	if err := compiler.compile(compiler.source.Block); err != nil {
		return nil, err
	}

	return compiler.instructions, nil
}
