package vm

type DataStack struct {
	elements []interface{}
}

func NewDataStack() *DataStack {
	return &DataStack{elements: []interface{}{}}
}

func (dataStack *DataStack) Push(elem interface{}) {
	dataStack.elements = append(dataStack.elements, elem)
}

func (dataStack *DataStack) Pop() interface{} {
	if len(dataStack.elements) == 0 {
		panic("data stack is empty")
	}

	elem := dataStack.Elem()
	dataStack.elements = dataStack.elements[:len(dataStack.elements)-1]
	return elem
}

func (dataStack *DataStack) Elem() interface{} {
	return dataStack.elements[len(dataStack.elements)-1]
}

