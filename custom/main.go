package main

import (
	"bytes"
	"fmt"
	"os"
)

// ./parse.go:179:3: undefined: sdfsf

func main() {
	file, err := os.Open("test.txt")
	if err != nil { panic(err) }
	defer file.Close()

	tokens := make([]Token, 0)
	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()

		tokens = append(tokens, Token{pos, tok, lit})
		// fmt.Printf("./%s:%d:%d\t%s\t%s\n", "test.txt", pos.line, pos.column, tok, lit)
		if tok == EOF {
			break
		}
	}

	tokenList := &Tokens{list: tokens}
	parser := Parser{tokenList}

	nodes := parser.ParseFile("input_test", tokenList) // TODO - token to represent file start?

	// WalkNode(nodes)


	buf := &genBuf{
		buf: new(bytes.Buffer),
	}
	GenerateCode(buf, nodes)
	fmt.Println(buf.String())

	err = os.WriteFile("./out/main.c", buf.buf.Bytes(), 0644)
	if err != nil { panic(err) }

	// s := bufio.NewScanner(file)
	// scanTokens(s)

	// dat, err := io.ReadAll(file)
	// if err != nil { panic(err) }

	// str := string(dat)
	// for _, b := range str {
	// 	fmt.Println(b)
	// }
}
