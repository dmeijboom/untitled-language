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


type OpenSection struct {
	Location *tokens.Location
}

func (openSection *OpenSection) Loc() *tokens.Location {
	return openSection.Location
}


type CloseSection struct {
	Location *tokens.Location
}

func (closeSection *CloseSection) Loc() *tokens.Location {
	return closeSection.Location
}


type MakeObject struct {
	Fields int
	Location *tokens.Location
}

func (makeObject *MakeObject) Loc() *tokens.Location {
	return makeObject.Location
}


type LoadType struct {
	Type TypeId
	Array bool
	Optional bool
	Location *tokens.Location
}

func (loadType *LoadType) Loc() *tokens.Location {
	return loadType.Location
}


type MakeType struct {
	Location *tokens.Location
}

func (makeType *MakeType) Loc() *tokens.Location {
	return makeType.Location
}


type MakeField struct {
	Location *tokens.Location
}

func (makeType *MakeField) Loc() *tokens.Location {
	return makeType.Location
}


type SetField struct {
	Location *tokens.Location
}

func (setField *SetField) Loc() *tokens.Location {
	return setField.Location
}


type NewObject struct {
	Fields int
	Location *tokens.Location
}

func (newObject *NewObject) Loc() *tokens.Location {
	return newObject.Location
}


type Initialize struct {
	Location *tokens.Location
}

func (init *Initialize) Loc() *tokens.Location {
	return init.Location
}


type LoadConst struct {
	Value interface{}
	Location *tokens.Location
}

func (loadConst *LoadConst) Loc() *tokens.Location {
	return loadConst.Location
}


type LoadVal struct {
 	Location *tokens.Location
}

func (loadVal *LoadVal) Loc() *tokens.Location {
	return loadVal.Location
}


type MakeCall struct {
	Args int
 	Location *tokens.Location
}

func (makeCall *MakeCall) Loc() *tokens.Location {
	return makeCall.Location
}


type LoadName struct {
	Name string
	Location *tokens.Location
}

func (loadName *LoadName) Loc() *tokens.Location {
	return loadName.Location
}


type StoreVal struct {
	HasValue bool
	Location *tokens.Location
}

func (storeVal *StoreVal) Loc() *tokens.Location {
	return storeVal.Location
}


/**
 * Instruction definitions
 */
func (setSection *OpenSection) instruction() {}
func (setSection *CloseSection) instruction() {}
func (makeField *MakeField) instruction() {}
func (loadType *LoadType) instruction() {}
func (loadName *LoadName) instruction() {}
func (makeType *MakeType) instruction() {}
func (storeVal *StoreVal) instruction() {}
func (loadConst *LoadConst) instruction() {}
func (loadVal *LoadVal) instruction() {}
func (setField *SetField) instruction() {}
func (newObject *NewObject) instruction() {}
func (makeObject *MakeObject) instruction() {}
func (makeCall *MakeCall) instruction() {}
func (initialize *Initialize) instruction() {}
