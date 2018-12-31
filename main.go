package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
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

	fmt.Println("LEXER")

	for _, token := range tokens {
		fmt.Println(token)
	}

	fmt.Println("\nPARSER")

	parser := NewParser(tokens)
	nodes, err := parser.Parse()

	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		data, _ := json.MarshalIndent(node, "", "  ")
		fmt.Println(string(data))
	}
}
