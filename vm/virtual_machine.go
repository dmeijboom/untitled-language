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

func (vm *VirtualMachine) processOpenSection(instruction *compiler.OpenSection) error {
	vm.callStack.Push(NewFrame(BlockFrame, instruction.Location))
	return nil
}

func (vm *VirtualMachine) processCloseSection(instruction *compiler.CloseSection) error {
	vm.callStack.Pop()
	return nil
}

func (vm *VirtualMachine) processMakeField(instruction *compiler.MakeField) error {
	fieldType := vm.dataStack.Pop().(*Type)
    vm.dataStack.Push(&ObjectField{
        Name: instruction.Name,
        Type: fieldType,
	})
	return nil
}

func (vm *VirtualMachine) processLoadType(instruction *compiler.LoadType) error {
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
        return nil
    }
    
    vm.dataStack.Push(&Type{
        Id: ArrayType,
        GenericParams: []Type{*rtype},
    })
	return nil
}

func (vm *VirtualMachine) processMakeType(instruction *compiler.MakeType) error {
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
	return nil
}

func (vm *VirtualMachine) processStoreVal(instruction *compiler.StoreVal) error {
	value := &Value{Mutable: true}

    if instruction.HasValue {
        value.Value = vm.dataStack.Pop().(interface{})
    }

    value.Type = vm.dataStack.Pop().(*Type)
    vm.callStack.Frame().Data[instruction.Name] = value

    println("@TODO: var validation")
	return nil
}

func (vm *VirtualMachine) processLoadConst(instruction *compiler.LoadConst) error {
	vm.dataStack.Push(instruction.Value)
	return nil
}

// func (vm *VirtualMachine) processLoadVal(instruction *compiler.LoadVal) error {
// 	return nil
// }

func (vm *VirtualMachine) processSetField(instruction *compiler.SetField) error {
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
	return nil
}

func (vm *VirtualMachine) processNewObject(instruction *compiler.NewObject) error {
	vm.dataStack.Push(NewObject())
	return nil
}

func (vm *VirtualMachine) processMakeObject(instruction *compiler.MakeObject) error {
	fields := []ObjectField{}

    for i := 0; i < instruction.Fields; i++ {
        fields = append(fields, *vm.dataStack.Pop().(*ObjectField))
    }

    vm.dataStack.Push(&ObjectDef{Fields: fields})
	return nil
}

func (vm *VirtualMachine) processInitialize(instruction *compiler.Initialize) error {
	println("@TODO: object validation")
	return nil
}


func (vm *VirtualMachine) Run() error {
	for vm.hasInstructions() {
		var err error

		switch instruction := vm.next().(type) {
		case *compiler.OpenSection:
			err = vm.processOpenSection(instruction)
			break
		case *compiler.CloseSection:
			err = vm.processCloseSection(instruction)
			break
		case *compiler.MakeField:
			err = vm.processMakeField(instruction)
			break
		case *compiler.LoadType:
			err = vm.processLoadType(instruction)
			break
		case *compiler.MakeType:
			err = vm.processMakeType(instruction)
			break
		case *compiler.StoreVal:
			err = vm.processStoreVal(instruction)
			break
		case *compiler.LoadConst:
			err = vm.processLoadConst(instruction)
			break
		// case *compiler.LoadVal:
		// 	err = vm.processLoadVal(instruction)
		// 	break
		case *compiler.SetField:
			err = vm.processSetField(instruction)
			break
		case *compiler.NewObject:
			err = vm.processNewObject(instruction)
			break
		case *compiler.MakeObject:
			err = vm.processMakeObject(instruction)
			break
		case *compiler.Initialize:
			err = vm.processInitialize(instruction)
			break
		default:
			panic(instruction)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
