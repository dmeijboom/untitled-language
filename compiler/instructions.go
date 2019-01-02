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


type MakeObject struct {
	Fields int
	Location *tokens.Location
}

func (makeObject *MakeObject) Loc() *tokens.Location {
	return makeObject.Location
}


type LoadType struct {
	TypeName string
	Type TypeId
	Array bool
	Optional bool
	Location *tokens.Location
}

func (loadType *LoadType) Loc() *tokens.Location {
	return loadType.Location
}


type MakeType struct {
	Name string
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
func (makeField *MakeField) instruction() {}
func (makeObject *MakeObject) instruction() {}
func (loadType *LoadType) instruction() {}
func (makeType *MakeType) instruction() {}
func (storeVal *StoreVal) instruction() {}
func (loadVal *LoadVal) instruction() {}
