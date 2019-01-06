package vm

type Array struct {
	values []*Value
}

func NewArray() *Array {
	return &Array{
		values: []*Value{},
	}
}

func (array *Array) Add(value *Value) {
	array.values = append(array.values, value)
}
