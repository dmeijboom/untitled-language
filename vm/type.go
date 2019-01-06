package vm

type TypeId int

const (
	ArrayType TypeId = iota
	StringType
	IntegerType
	FloatType
	BooleanType
	FunctionType
	ObjectType
)

type Type struct {
	Id TypeId
	Name string
	ObjectDef *ObjectDef
	GenericParams []Type
}

func (type_ *Type) Compare(otherType *Type) bool {
	if type_.Id != otherType.Id ||
		type_.Name != otherType.Name ||
		len(type_.GenericParams) != len(otherType.GenericParams) {
		return false
	}

	if len(type_.GenericParams) > 0 {
		for i := 0; i < len(type_.GenericParams); i++ {
			if !(&type_.GenericParams[i]).Compare(&otherType.GenericParams[i]) {
				return false
			}
		}
	}

	return true
}
