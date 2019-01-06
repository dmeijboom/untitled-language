package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"encoding/json"
	"dmeijboom/config/vm"
	"dmeijboom/config/compiler"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	content, err := ioutil.ReadFile("./config/filesystem.cf")

	if err != nil {
		panic(err)
	}

	lexer := NewLexer(string(content))
	tokens, err := lexer.Lex()

	if err != nil {
		panic(err)
	}

	fmt.Println("LEXER\n---")

	for _, token := range tokens {
		fmt.Println(token)
	}

	fmt.Println("\nPARSER\n---")

	parser := NewParser(tokens)
	source, err := parser.Parse()

	if err != nil {
		panic(err)
	}

	spew.Config.DisablePointerAddresses = true
	spew.Dump(source.Block.Body)

	fmt.Println("\nCOMPILER\n---")

	compiler := compiler.NewCompiler(source)
	instructions, err := compiler.Compile()

	if err != nil {
		panic(err)
	}

	for _, instruction := range instructions {
		data, _ := json.Marshal(instruction)
		fmt.Println(reflect.TypeOf(instruction).Elem().Name() + " " + string(data))
	}

	fmt.Println("\nRUN\n---")

	machine := vm.NewVm(instructions)
	machine.Set("writeln", &vm.Value{
		Type: &vm.Type{Id: vm.FunctionType},
		Mutable: false,
		Value: &vm.Function{
			Name: "writeln",
			Func: func(values []*vm.Value) {
				fmt.Println(values[0].Value)
			},
		},
	})
	machine.Set("array_add", &vm.Value{
		Type: &vm.Type{Id: vm.FunctionType},
		Mutable: false,
		Value: &vm.Function{
			Name: "add",
			Func: func(values []*vm.Value) {
				values[0].Value.(*vm.Array).Add(values[1])
			},
		},
	})

	err = machine.Run()

	if err != nil {
		panic(err)
	}
}
