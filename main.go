package main

import (
	"fmt"
	"io/ioutil"
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

	for _, token := range tokens {
		fmt.Println(token)
	}
}
