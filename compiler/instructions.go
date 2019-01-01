package compiler

import (
	"dmeijboom/config/tokens"
)

type TypeId int

const (
	StringType TypeId = iota
	IntegerType
	BooleanType
	FloatType
	ObjectType
	UserType
)

type Instruction interface {
	instruction()
	Loc() *tokens.Location
}


type SetSection struct {
	Name string
	Location *tokens.Location
}

func (setSection *SetSection) Loc() *tokens.Location {
	return setSection.Location
}


type MakeType struct {
	TypeName string
	Type TypeId
	Fields int
	Array bool
	Optional bool
	Location *tokens.Location
}

func (makeType *MakeType) Loc() *tokens.Location {
	return makeType.Location
}


type MakeField struct {
	Name string
	Location *tokens.Location
}

func (makeType *MakeField) Loc() *tokens.Location {
	return makeType.Location
}


type LoadVal struct {
	Value interface{}
	Location *tokens.Location
}

func (loadVal *LoadVal) Loc() *tokens.Location {
	return loadVal.Location
}


type StoreVal struct {
	Name string
	HasValue bool
	Location *tokens.Location
}

func (storeVal *StoreVal) Loc() *tokens.Location {
	return storeVal.Location
}


/**
 * Instruction definitions
 */
func (setSection *SetSection) instruction() {}
func (makeType *MakeType) instruction() {}
func (makeField *MakeField) instruction() {}
func (storeVal *StoreVal) instruction() {}
func (loadVal *LoadVal) instruction() {}
