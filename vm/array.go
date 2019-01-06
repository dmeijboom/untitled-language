package vm

type Array struct {
	data []interface{}
}

func NewArray() *Array {
	return &Array{
		data: []interface{}{},
	}
}
