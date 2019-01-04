package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"encoding/json"
	"dmeijboom/config/vm"
	"dmeijboom/config/compiler"
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

	for _, node := range source.Block.Body {
		data, _ := json.MarshalIndent(node, "", "  ")
		fmt.Println(reflect.TypeOf(node).Elem().Name() + " " + string(data))
	}

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

	vm := vm.NewVm(instructions)
	err = vm.Run()

	if err != nil {
		panic(err)
	}
}
