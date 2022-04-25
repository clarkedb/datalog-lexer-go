package main

import "fmt"

type Token struct {
	TokenType  string
	Value      string
	LineNumber int
}

func (t Token) String() string {
	// Should be (TYPE,"value",1)
	return fmt.Sprintf("(%v,\"%s\",%d)", t.TokenType, t.Value, t.LineNumber)
}
