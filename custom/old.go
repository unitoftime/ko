package main

// func scanTokens(scanner *bufio.Scanner) []Token {
// 	// file := "test.txt"
// 	line := 1

// 	ret := make([]Token, 0, 2048) // TODO:
// 	for scanner.Scan() {
// 		lineStr := scanner.Text()
// 		for charCol, r := range lineStr {
// 			pos := Pos{line: line, char: charCol}
// 			tok := Token{
// 				Pos: pos,
// 				Type: getTokenType(pos, r),
// 				Lex: "???",
// 			}

// 			ret = append(ret, tok)
// 		}
// 		fmt.Println(scanner.Text())
// 	}

// 	err := scanner.Err()
// 	if err != nil { panic(err) }

// 	return ret
// }

// func getTokenType(pos Pos, r rune) TokenType {
// 	switch r {
// 	case '(':
// 		return TokenLParen
// 	case ')':
// 		return TokenRParen
// 	case '{':
// 		return TokenLBrace
// 	case '}':
// 		return TokenRBrace
// 	case ',':
// 		return TokenComma
// 	case '.':
// 		return TokenDot
// 	case '-':
// 		return TokenMinus
// 	case '+':
// 		return TokenPlus
// 	case ';':
// 		return TokenSemicolon
// 	case '/':
// 		return TokenSlash
// 	case '*':
// 		return TokenStar
// 	}
// 	panic(fmt.Sprintf("UnknownToken: %c", r))
// }

// type Pos struct {
// 	// file string // TODO:???
// 	line, char int
// }

// type Token struct {
// 	Type TokenType
// 	Lex string
// 	Pos Pos
// }

// type TokenType uint8
// const (
// 	// TokenInvalid

// 	//Single-char
// 	TokenLParen TokenType = iota
// 	TokenRParen
// 	TokenLBrace
// 	TokenRBrace
// 	TokenComma
// 	TokenDot
// 	TokenMinus
// 	TokenPlus
// 	TokenSemicolon
// 	TokenSlash
// 	TokenStar

// 	// One or two character
// 	TokenBang
// 	TokenBangEqual
// 	TokenEqual
// 	TokenDoubleEqual
// 	TokenGreater
// 	TokenGreaterEqual
// 	TokenLess
// 	TokenLessEqual

// 	//Literals
//   // IDENTIFIER, STRING, NUMBER,

//   // // Keywords.
//   // AND, CLASS, ELSE, FALSE, FUN, FOR, IF, NIL, OR,
//   // PRINT, RETURN, SUPER, THIS, TRUE, VAR, WHILE,

//   // EOF
// )
