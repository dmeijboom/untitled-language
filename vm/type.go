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
	Optional bool
	ObjectDef *ObjectDef
	GenericParams []Type
}

func (type_ *Type) FullName() string {
	typeName := type_.Name

	if type_.Id == ArrayType {
		typeName = "[]" + (&type_.GenericParams[0]).FullName()
	} else if len(type_.GenericParams) > 0 {
		typeName += "["

		for i := 0; i < len(type_.GenericParams); i++ {
			typeName += (&type_.GenericParams[i]).FullName()

			if i < len(type_.GenericParams)-1 {
				typeName += ", "
			}
		}

		typeName += "]"
	}

	if type_.Optional {
		typeName += "?"
	}

	return typeName
}

func (type_ *Type) Equals(otherType *Type) bool {
	if type_.Id != otherType.Id ||
		type_.Name != otherType.Name ||
		len(type_.GenericParams) != len(otherType.GenericParams) {
		return false
	}

	if len(type_.GenericParams) > 0 {
		for i := 0; i < len(type_.GenericParams); i++ {
			if !(&type_.GenericParams[i]).Equals(&otherType.GenericParams[i]) {
				return false
			}
		}
	}

	return true
}
