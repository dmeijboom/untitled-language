package vm

// type FunctionArgument struct {
// 	Name string
// 	Type *Type
// }

type GoFunc func (values []*Value)

type Function struct {
	Name string
	Func GoFunc
	// Args []FunctionArgument
}
