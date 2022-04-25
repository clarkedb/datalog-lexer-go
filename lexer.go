package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Lexer struct {
	Filename    string
	CurrentLine int
	Tokens      []Token
}

func (l Lexer) String() string {
	var tokenStrings []string
	for _, t := range l.Tokens {
		tokenStrings = append(tokenStrings, fmt.Sprintf("%v", t))
	}
	return fmt.Sprintf("%s\nTotal Tokens = %d", strings.Join(tokenStrings, "\n"), len(l.Tokens))
}

func (l Lexer) Length() int {
	return len(l.Tokens)
}

func (l *Lexer) Tokenize() {

	f, err := os.Open(l.Filename)
	if err != nil {
		panic(err)
	}

	var expr []rune
	var c rune
	var n rune
	var keyword bool

	r := bufio.NewReader(f)

	for {
		_, err := peek(r)
		if err == io.EOF {
			break
		}
		expr = nil
		c = get(r)
		if unicode.IsSpace(c) {
			// NEW LINE
			if c == '\n' {
				l.CurrentLine += 1
			}
		} else if unicode.IsLetter(c) {
			keyword = false
			n, _ = peek(r)

			expr = append(expr, c)
			if !isAlphaNumeric(n) {
				l.addToken("ID", expr, l.CurrentLine)
			} else {
				if c == 'S' {
					// SCHEMES
					if n == 'c' {
						c = get(r)
						expr = append(expr, c)
						n, _ = peek(r)

						if n == 'h' {
							c = get(r)
							expr = append(expr, c)
							n, _ = peek(r)

							if n == 'e' {
								c = get(r)
								expr = append(expr, c)
								n, _ = peek(r)

								if n == 'm' {
									c = get(r)
									expr = append(expr, c)
									n, _ = peek(r)

									if n == 'e' {
										c = get(r)
										expr = append(expr, c)
										n, _ = peek(r)

										if n == 's' {
											c = get(r)
											expr = append(expr, c)
											n, _ = peek(r)

											if !isAlphaNumeric(n) {
												l.addToken("SCHEMES", expr, l.CurrentLine)
												keyword = true
											}
										}
									}
								}
							}
						}
					}
				} else if c == 'F' {
					// FACTS
					if n == 'a' {
						c = get(r)
						expr = append(expr, c)
						n, _ = peek(r)

						if n == 'c' {
							c = get(r)
							expr = append(expr, c)
							n, _ = peek(r)

							if n == 't' {
								c = get(r)
								expr = append(expr, c)
								n, _ = peek(r)

								if n == 's' {
									c = get(r)
									expr = append(expr, c)
									n, _ = peek(r)

									if !isAlphaNumeric(n) {
										l.addToken("FACTS", expr, l.CurrentLine)
										keyword = true
									}
								}
							}
						}
					}
				} else if c == 'R' {
					// RULES
					if n == 'u' {
						c = get(r)
						expr = append(expr, c)
						n, _ = peek(r)

						if n == 'l' {
							c = get(r)
							expr = append(expr, c)
							n, _ = peek(r)

							if n == 'e' {
								c = get(r)
								expr = append(expr, c)
								n, _ = peek(r)

								if n == 's' {
									c = get(r)
									expr = append(expr, c)
									n, _ = peek(r)

									if !isAlphaNumeric(n) {
										l.addToken("RULES", expr, l.CurrentLine)
										keyword = true
									}
								}
							}
						}
					}
				} else if c == 'Q' {
					// QUERIES
					if n == 'u' {
						c = get(r)
						expr = append(expr, c)
						n, _ = peek(r)

						if n == 'e' {
							c = get(r)
							expr = append(expr, c)
							n, _ = peek(r)

							if n == 'r' {
								c = get(r)
								expr = append(expr, c)
								n, _ = peek(r)

								if n == 'i' {
									c = get(r)
									expr = append(expr, c)
									n, _ = peek(r)

									if n == 'e' {
										c = get(r)
										expr = append(expr, c)
										n, _ = peek(r)

										if n == 's' {
											c = get(r)
											expr = append(expr, c)
											n, _ = peek(r)

											if !isAlphaNumeric(n) {
												l.addToken("QUERIES", expr, l.CurrentLine)
												keyword = true
											}
										}
									}
								}
							}
						}
					}
				}

				if !keyword {
					for !keyword && isAlphaNumeric(n) {
						c = get(r)
						expr = append(expr, c)
						n, _ = peek(r)
					}
					l.addToken("ID", expr, l.CurrentLine)
				}
			}
		} else if isLegalOpChar(c) {
			// OPERATION
			expr = append(expr, c)
			switch c {
			case ':':
				n, _ = peek(r)
				if n == '-' {
					c = get(r)
					expr = append(expr, c)
					l.addToken("COLON_DASH", expr, l.CurrentLine)
				} else {
					l.addToken("COLON", expr, l.CurrentLine)
				}
			case ',':
				l.addToken("COMMA", expr, l.CurrentLine)
			case '.':
				l.addToken("PERIOD", expr, l.CurrentLine)
			case '?':
				l.addToken("Q_MARK", expr, l.CurrentLine)
			case '(':
				l.addToken("LEFT_PAREN", expr, l.CurrentLine)
			case ')':
				l.addToken("RIGHT_PAREN", expr, l.CurrentLine)
			case '*':
				l.addToken("MULTIPLY", expr, l.CurrentLine)
			case '+':
				l.addToken("ADD", expr, l.CurrentLine)
			}
		} else if c == '#' {
			// COMMENT
			expr = append(expr, c)
			n, _ = peek(r)

			if n == '|' {
				// block comment case
				c = get(r)
				expr = append(expr, c)
				n, err = peek(r)

				commentStartLine := l.CurrentLine
				validComment := true
				endComment := false

				for !endComment && validComment {
					if err == io.EOF {
						l.addToken("UNDEFINED", expr, commentStartLine)
						validComment = false
					} else {
						c = get(r)
						expr = append(expr, c)
						if c == '\n' {
							l.CurrentLine += 1
						}
						n, err = peek(r)
					}

					if c == '|' && n == '#' {
						endComment = true
						c = get(r)
						expr = append(expr, c)
					}
				}

				if validComment {
					l.addToken("COMMENT", expr, commentStartLine)
				}
			} else {
				// single line comment case
				for n != '\n' && err != io.EOF {
					c = get(r)
					expr = append(expr, c)
					n, err = peek(r)
				}
				l.addToken("COMMENT", expr, l.CurrentLine)
			}

		} else if c == '\'' {
			// STRING
			expr = append(expr, c)
			n, err = peek(r)

			numQuotes := 1
			stringStartLine := l.CurrentLine
			validString := true
			endString := false

			for !endString && validString {
				if err == io.EOF {
					l.addToken("UNDEFINED", expr, stringStartLine)
					validString = false
				} else {
					c = get(r)
					expr = append(expr, c)
					switch c {
					case '\'':
						numQuotes += 1
					case '\n':
						l.CurrentLine += 1
					}
					n, err = peek(r)
				}

				if (numQuotes%2 == 0) && (c == '\'') && (n != '\'') {
					endString = true
				}
			}

			if validString {
				l.addToken("STRING", expr, stringStartLine)
			}
		} else {
			// UNDEFINED
			expr = append(expr, c)
			l.addToken("UNDEFINED", expr, l.CurrentLine)
		}
	}

	// Final token is EOF every time
	l.addToken("EOF", []rune{}, l.CurrentLine)

	f.Close()
}

func (l *Lexer) addToken(tokenType string, valueList []rune, lineNumber int) {
	// value := strings.Join(valueList, "")
	value := string(valueList)
	l.Tokens = append(l.Tokens, Token{TokenType: tokenType, Value: value, LineNumber: lineNumber})
}

func isLegalOpChar(c rune) bool {
	return strings.Contains(",.?():*+", string(c))
}

func isAlphaNumeric(c rune) bool {
	return (unicode.IsNumber(c) || unicode.IsLetter(c))
}

func peek(r *bufio.Reader) (rune, error) {
	bytes, err := r.Peek(1)
	if err == io.EOF {
		return ' ', io.EOF
	} else if err != nil {
		panic(err)
	}
	return rune(bytes[0]), nil
}

func get(r *bufio.Reader) rune {
	c, _, err := r.ReadRune()
	if err != nil {
		panic(err)
	}
	return c
}
