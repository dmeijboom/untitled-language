package vm

// type FunctionArgument struct {
// 	Name string
// 	Type *Type
// }

type GoFunc func (values []*Value)

type FunctionLookup struct {
	Name string
	Value *Value
}

type Function struct {
	Name string
	Func GoFunc
	// Args []FunctionArgument
}
