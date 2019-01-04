package vm

type ObjectField struct {
	Name string
	Type *Type
}

type ObjectDef struct {
	Fields []ObjectField
}

func (objectDef *ObjectDef) FieldByName(name string) *ObjectField {
	for _, field := range objectDef.Fields {
		if field.Name == name {
			return &field
		}
	}

	return nil
}


type Object struct {
	Fields map[string]*Value
}

func NewObject() *Object {
	return &Object{
		Fields: map[string]*Value{},
	}
}
