package vm

type Value struct {
	Type *Type
	Mutable bool
	Value interface{}
}
