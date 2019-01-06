package vm

import (
	"fmt"
	"dmeijboom/config/compiler"
)

type VirtualMachine struct {
	index int
	root *Frame
	callStack *CallStack
	dataStack *DataStack
	instructions []compiler.Instruction
}

func NewVm(instructions []compiler.Instruction) *VirtualMachine {
	vm := &VirtualMachine{
		root: NewFrame(RootFrame, nil),
		callStack: NewCallStack(),
		dataStack: NewDataStack(),
		instructions: instructions,
	}

	vm.callStack.Push(vm.root)

	return vm
}

func (vm *VirtualMachine) Set(name string, value *Value) {
	vm.root.Data[name] = value
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
	vm.dataStack.Pop()
	vm.callStack.Push(NewFrame(BlockFrame, instruction.Location))
	return nil
}

func (vm *VirtualMachine) processCloseSection(instruction *compiler.CloseSection) error {
	vm.callStack.Pop()
	return nil
}

func (vm *VirtualMachine) processMakeField(instruction *compiler.MakeField) error {
	name := vm.dataStack.Pop().(string)
	fieldType := vm.dataStack.Pop().(*Type)

    vm.dataStack.Push(&ObjectField{
        Name: name,
        Type: fieldType,
	})
	return nil
}

func (vm *VirtualMachine) processLoadType(instruction *compiler.LoadType) error {
	var rtype *Type

	typeName := vm.dataStack.Pop().(string)

    if instruction.Type == compiler.UserType {
        rtype = vm.lookupType(typeName)

        if rtype == nil {
            panic("Type not found: " + typeName)
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
	name := vm.dataStack.Pop().(string)
	def := vm.dataStack.Pop()

    if objectDef, ok := def.(*ObjectDef); ok {
        vm.callStack.Frame().Types[name] = &Type{
            Id: ObjectType,
            Name: name,
            ObjectDef: objectDef,
        }
    } else {
        panic("Not supported")
	}
	
	return nil
}

func (vm *VirtualMachine) processStoreVal(instruction *compiler.StoreVal) error {
	name := vm.dataStack.Pop().(string)
	value := &Value{
		Mutable: true,
	}

    if instruction.HasValue {
        value.Value = vm.dataStack.Pop().(interface{})
    }

	value.Type = vm.dataStack.Pop().(*Type)

    vm.callStack.Frame().Data[name] = value

    // println("@TODO: var validation")
	return nil
}

func (vm *VirtualMachine) processLoadConst(instruction *compiler.LoadConst) error {
	vm.dataStack.Push(instruction.Value)
	return nil
}

func (vm *VirtualMachine) processLoadName(instruction *compiler.LoadName) error {
	vm.dataStack.Push(instruction.Name)
	return nil
}

func (vm *VirtualMachine) processLoadMember(instruction *compiler.LoadMember) error {
	name := vm.dataStack.Pop().(string)
	value := vm.dataStack.Pop().(*Value)

	if value.Type.Id != ObjectType {
		return fmt.Errorf("Cannot use non-object type as object")
	}

	object := value.Value.(*Object)

	if field, exist := object.Fields[name]; exist {
		vm.dataStack.Push(field)
		return nil
	}

	return fmt.Errorf("%s does not contain the `%s` field", value.Type.Name, name)
}

func (vm *VirtualMachine) processLoadVal(instruction *compiler.LoadVal) error {
	name := vm.dataStack.Pop().(string)

	if value := vm.callStack.Frame().Get(name); value != nil {
		vm.dataStack.Push(value)
		return nil
	}

	return fmt.Errorf("Name `%s` not found", name)
}

func (vm *VirtualMachine) processSetField(instruction *compiler.SetField) error {
	value := vm.dataStack.Pop().(interface{})
	fieldName := vm.dataStack.Pop().(string)
	object := vm.dataStack.Pop().(*Object)
	objectType := vm.dataStack.Elem().(*Type)
    field := objectType.ObjectDef.FieldByName(fieldName)

    object.Fields[fieldName] = &Value{
        Type: field.Type,
        Mutable: true,
        Value: value,
    }

    vm.dataStack.Push(object)

	// println("@TODO: field validation")
	return nil
}

func (vm *VirtualMachine) processNewObject(instruction *compiler.NewObject) error {
	vm.dataStack.Push(NewObject())
	return nil
}

func (vm *VirtualMachine) processMakeCall(instruction *compiler.MakeCall) error {
	callable := vm.dataStack.Pop().(*Value)

	if callable.Type.Id != FunctionType {
		return fmt.Errorf("Cannot call non-function type")
	}

	fn := callable.Value.(*Function)
	args := []*Value{}

	for i := 0; i < instruction.Args; i++ {
		args = append(args, vm.dataStack.Pop().(*Value))
	}

	fn.Func(args)

	return nil
}

func (vm *VirtualMachine) processMakeObject(instruction *compiler.MakeObject) error {
	vm.dataStack.Pop() // unused 'object' typename

	fields := []ObjectField{}

    for i := 0; i < instruction.Fields; i++ {
        fields = append(fields, *vm.dataStack.Pop().(*ObjectField))
    }

    vm.dataStack.Push(&ObjectDef{Fields: fields})
	return nil
}

func (vm *VirtualMachine) processInitialize(instruction *compiler.Initialize) error {
	// println("@TODO: object validation")
	return nil
}


func (vm *VirtualMachine) Run() error {
	for vm.hasInstructions() {
		var err error
		instr := vm.next()

		switch instruction := instr.(type) {
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
		case *compiler.MakeCall:
			err = vm.processMakeCall(instruction)
			break
		case *compiler.StoreVal:
			err = vm.processStoreVal(instruction)
			break
		case *compiler.LoadConst:
			err = vm.processLoadConst(instruction)
			break
		case *compiler.LoadName:
			err = vm.processLoadName(instruction)
			break
		case *compiler.LoadVal:
			err = vm.processLoadVal(instruction)
			break
		case *compiler.LoadMember:
			err = vm.processLoadMember(instruction)
			break
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
			return fmt.Errorf("%s at %d:%d", err.Error(), instr.Loc().Line, instr.Loc().Column)
		}
	}

	return nil
}
