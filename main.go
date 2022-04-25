package main

import (
	"fmt"
	"os"
)

func main() {
	lexer := new(Lexer)
	lexer.Filename = os.Args[1]
	lexer.CurrentLine = 1
	lexer.Tokenize()
	fmt.Println(lexer)
}
