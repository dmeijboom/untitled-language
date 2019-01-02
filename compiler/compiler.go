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
	case *ast.Initialize:
		break
	case *ast.Block:
		for _, subNode := range node.Body {
			if err := compiler.compile(subNode); err != nil {
				return err
			}
		}
		break
	case *ast.Literal:
		compiler.add(&LoadVal{
			Value: node.Value,
			Location: node.Loc(),
		})
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
		if node.Name.Value == "object" {
			if len(node.Fields) > 0 {
				for _, field := range node.Fields {
					if err := compiler.compile(&field); err != nil {
						return err
					}
				}
			}

			compiler.add(&MakeObject{
				Fields: len(node.Fields),
				Location: node.Loc(),
			})
			break
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
		break
	case *ast.Typedef:
		if err := compiler.compile(node.Type); err != nil {
			return err
		}

		compiler.add(&MakeType{
			Name: node.Name.Value,
			Location: node.Loc(),
		})
		break
	case *ast.Assign:
		if err := compiler.compile(node.Type); err != nil {
			return err
		}

		if node.Value != nil {
			if err := compiler.compile(node.Value); err != nil {
				return err
			}
		}

		compiler.add(&StoreVal{
			Name: node.Name.Value,
			HasValue: node.Value != nil,
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
