package vm

import (
	"dmeijboom/config/compiler"
)

type VirtualMachine struct {
	index int
	callStack *CallStack
	dataStack *DataStack
	instructions []compiler.Instruction
}

func NewVm(instructions []compiler.Instruction) *VirtualMachine {
	return &VirtualMachine{
		callStack: NewCallStack(),
		dataStack: NewDataStack(),
		instructions: instructions,
	}
}

func (vm *VirtualMachine) hasInstructions() bool {
	return vm.index <= len(vm.instructions)-1
}

func (vm *VirtualMachine) next() compiler.Instruction {
	instr := vm.instructions[vm.index]
	vm.index++
	return instr
}

func (vm *VirtualMachine) peek() compiler.Instruction {
	return vm.instructions[vm.index]
}

func (vm *VirtualMachine) convertType(compilerType compiler.TypeId) *Type {
	switch compilerType {
	case compiler.StringType:
		return &Type{Id: StringType}
	case compiler.IntegerType:
		return &Type{Id: IntegerType}
	case compiler.BooleanType:
		return &Type{Id: BooleanType}
	case compiler.FloatType:
		return &Type{Id: FloatType}
	}

	panic("Unknown type in convertType")
}

func (vm *VirtualMachine) lookupType(name string) *Type {
	frame := vm.callStack.Frame()

	for frame != nil {
		if type_, found := frame.Types[name]; found {
			return type_
		}

		frame = frame.Parent
	}

	return nil
}

func (vm *VirtualMachine) Run() error {
	for vm.hasInstructions() {
		instr := vm.peek()

		switch instruction := instr.(type) {
		case *compiler.SetSection:
			vm.callStack.Push(NewFrame(BlockFrame, instruction.Location))
			break
		case *compiler.LoadConst:
			vm.dataStack.Push(instruction.Value)
			break
		case *compiler.LoadType:
			var rtype *Type

			if instruction.Type == compiler.UserType {
				rtype = vm.lookupType(instruction.TypeName)

				if rtype == nil {
					panic("Type not found: " + instruction.TypeName)
				}
			} else {
				rtype = vm.convertType(instruction.Type)
			}

			if !instruction.Array {
				vm.dataStack.Push(rtype)
				break
			}
			
			vm.dataStack.Push(&Type{
				Id: ArrayType,
				GenericParams: []Type{*rtype},
			})
			break
		case *compiler.MakeField:
			fieldType := vm.dataStack.Pop().(*Type)
			vm.dataStack.Push(&ObjectField{
				Name: instruction.Name,
				Type: fieldType,
			})
			break
		case *compiler.MakeObject:
			fields := []ObjectField{}

			for i := 0; i < instruction.Fields; i++ {
				fields = append(fields, *vm.dataStack.Pop().(*ObjectField))
			}

			vm.dataStack.Push(&ObjectDef{Fields: fields})
			break
		case *compiler.MakeType:
			def := vm.dataStack.Pop()

			if objectDef, ok := def.(*ObjectDef); ok {
				vm.callStack.Frame().Types[instruction.Name] = &Type{
					Id: ObjectType,
					Name: instruction.Name,
					ObjectDef: objectDef,
				}
			} else {
				panic("Not supported")
			}
			break
		case *compiler.StoreVal:
			value := &Value{Mutable: true}

			if instruction.HasValue {
				value.Value = vm.dataStack.Pop().(interface{})
			}

			value.Type = vm.dataStack.Pop().(*Type)
			vm.callStack.Frame().Data[instruction.Name] = value

			println("@TODO: var validation")
			break
		case *compiler.NewObject:
			vm.dataStack.Push(NewObject())
			break
		case *compiler.SetField:
			value := vm.dataStack.Pop().(interface{})
			object := vm.dataStack.Pop().(*Object)
			objectType := vm.dataStack.Elem().(*Type)
			field := objectType.ObjectDef.FieldByName(instruction.Name)

			object.Fields[instruction.Name] = &Value{
				Type: field.Type,
				Mutable: true,
				Value: value,
			}

			vm.dataStack.Push(object)

			println("@TODO: field validation")
			break
		case *compiler.Initialize:
			println("@TODO: object validation")
			break
		default:
			panic(instruction)
		}

		vm.next()
	}

	return nil
}
