package main

import (
	"bufio"
	"io"
	"unicode"
)

type Token struct {
	pos Position
	token TokenType
	str string
}

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL

	// Single Character
	SEMI // ;
	COMMA // ;
	DOT // .
	MUL // *
	DIV // /
	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }
	LBRACK // [
	RBRACK // ]

	// One or two character
	ADD // +
	INC // ++
	PLUSEQ // +=
	SUB // -
	DEC // --
	SUBEQ // -=

	BANG // !
	BANGEQUAL // !=
	EQUAL // =
	EQUALEQUAL // ==
	GREATER // >
	GREATEREQUAL // >=
	LESS // <
	LESSEQUAL // <=
	AND // &
	ANDAND // &&
	OR // ||
	COLON // :
	WALRUS // :=

	HASH // #

	// Literals
	IDENT
	INT
	FLOAT
	STRING
	CHARLIT
	LINECOMMENT

  // Keywords.
	PACKAGE
	FUNC
	STRUCT
	TYPE
	VAR
	CONST
	RETURN
	IF
	ELSE
	FOR
	TRUE
	FALSE
	FOREIGN
	NIL
	SWITCH
	CASE
	BREAK
	DEFAULT
	ENUM
  // PRINT, THIS
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",

	// Single character
	SEMI:    ";",
	COMMA:    ",",
	DOT: ".",
	MUL: "MUL",
	DIV: "DIV",
	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",
	LBRACK: "[",
	RBRACK: "]",
	HASH: "#",

	// One or two character
	ADD: "ADD",
	INC: "INC",
	PLUSEQ: "PLUSEQ",
	SUB: "SUB",
	DEC: "DEC",
	SUBEQ: "SUBEQ",

	BANG: "!",
	BANGEQUAL: "!=",
	EQUAL: "=",
	EQUALEQUAL: "==",
	GREATER: ">",
	GREATEREQUAL: ">=",
	LESS: "<",
	LESSEQUAL: "<=",
	AND: "&",
	ANDAND: "&&",
	OR: "||",
	COLON: ":",
	WALRUS: ":=",

	IDENT:   "IDENT",
	INT:     "INT",
	FLOAT:     "FLOAT",
	STRING:     "STRING",
	CHARLIT:     "CHARLIT",

	LINECOMMENT: "LINECOMMENT",


	// Keywords
	PACKAGE: "package",
	FUNC: "func",
	STRUCT: "struct",
	TYPE: "type",
	VAR: "var",
	CONST: "const",
	RETURN: "return",
	IF: "if",
	ELSE: "else",
	FOR: "for",
	TRUE: "true",
	FALSE: "false",
	FOREIGN: "foreign",
	NIL: "nil",
	SWITCH: "switch",
	CASE: "case",
	BREAK: "break",
	DEFAULT: "default",
	ENUM: "enum",
}
var keywordList = map[string]TokenType{
	tokens[PACKAGE]: PACKAGE,
	tokens[FUNC]: FUNC,
	tokens[STRUCT]: STRUCT,
	tokens[TYPE]: TYPE,
	tokens[VAR]: VAR,
	tokens[CONST]: CONST,
	tokens[RETURN]: RETURN,
	tokens[IF]: IF,
	tokens[ELSE]: ELSE,
	tokens[FOR]: FOR,
	tokens[TRUE]: TRUE,
	tokens[FALSE]: FALSE,
	tokens[FOREIGN]: FOREIGN,
	tokens[NIL]: NIL,
	tokens[SWITCH]: SWITCH,
	tokens[CASE]: CASE,
	tokens[BREAK]: BREAK,
	tokens[DEFAULT]: DEFAULT,
	tokens[ENUM]: ENUM,
}

func (t TokenType) String() string {
	return tokens[t]
}

type Position struct {
	filename string
	line   int
	column int
}

type Lexer struct {
	lastToken TokenType
	pos    Position
	reader *bufio.Reader
}

func NewLexer(filename string, reader io.Reader) *Lexer {
	return &Lexer{
		lastToken: ILLEGAL,
		pos:    Position{filename: filename, line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (l *Lexer) Lex() (Position, TokenType, string) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, "EOF"
			}

			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}

		// update the column to the position of the newly read in rune
		l.pos.column++

		switch r {
		case '\n':
			// Decide if we want to add semicolon
			if l.lastToken == IDENT || l.lastToken == RPAREN || l.lastToken == RBRACE || l.lastToken == RBRACK || l.lastToken == INT || l.lastToken == FLOAT || l.lastToken == INC || l.lastToken == DEC || l.lastToken == STRING || l.lastToken == CHARLIT || l.lastToken == TRUE || l.lastToken == FALSE {
				l.lastToken = SEMI
				l.resetPosition()
				return l.pos, SEMI, ";"
			}
			l.resetPosition()
		case ';':
			l.lastToken = SEMI
			return l.pos, SEMI, ";"
		case ',':
			l.lastToken = COMMA
			return l.pos, COMMA, ","
		case '.':
			l.lastToken = DOT
			return l.pos, DOT, ","
		case '*':
			l.lastToken = MUL
			return l.pos, MUL, "*"
		case '(':
			l.lastToken = LPAREN
			return l.pos, LPAREN, "("
		case ')':
			l.lastToken = RPAREN
			return l.pos, RPAREN, ")"
		case '{':
			l.lastToken = LBRACE
			return l.pos, LBRACE, "}"
		case '}':
			l.lastToken = RBRACE
			return l.pos, RBRACE, "}"
		case '[':
			l.lastToken = LBRACK
			return l.pos, LBRACK, "["
		case ']':
			l.lastToken = RBRACK
			return l.pos, RBRACK, "]"
		case '#':
			l.lastToken = HASH
			return l.pos, HASH, "#"

			//--------------------------------------------------------------------------------
			// - One or two tokens
			//--------------------------------------------------------------------------------
		case '+':
			if l.match('+') {
				l.lastToken = INC
				return l.pos, INC, "++"
			} else if l.match('=') {
				l.lastToken = PLUSEQ
				return l.pos, l.lastToken, "+="
			}

			l.lastToken = ADD
			return l.pos, ADD, "+"
		case '-':
			if l.match('-') {
				l.lastToken = DEC
				return l.pos, DEC, "--"
			} else if l.match('=') {
				l.lastToken = SUBEQ
				return l.pos, l.lastToken, "-="
			}

			l.lastToken = SUB
			return l.pos, SUB, "-"
		case '!':
			if l.match('=') {
				l.lastToken = BANGEQUAL
				return l.pos, BANGEQUAL, "!="
			}
			l.lastToken = BANG
			return l.pos, BANG, "!"
		case '=':
			if l.match('=') {
				l.lastToken = EQUALEQUAL
				return l.pos, EQUALEQUAL, "=="
			}
			l.lastToken = EQUAL
			return l.pos, EQUAL, "="
		case '>':
			if l.match('=') {
				l.lastToken = GREATEREQUAL
				return l.pos, GREATEREQUAL, ">="
			}
			l.lastToken = GREATER
			return l.pos, GREATER, ">"
		case '<':
			if l.match('=') {
				l.lastToken = LESSEQUAL
				return l.pos, LESSEQUAL, "<="
			}
			l.lastToken = LESS
			return l.pos, LESS, "<"
		case '&':
			if l.match('&') {
				l.lastToken = ANDAND
				return l.pos, l.lastToken, "&&"
			}
			l.lastToken = AND
			return l.pos, l.lastToken, "&"

		case '|':
			if l.match('|') {
				l.lastToken = OR
				return l.pos, l.lastToken, "||"
			}
			panic("bitwise OR not yet impl")
		case ':':
			if l.match('=') {
				l.lastToken = WALRUS
				return l.pos, l.lastToken, ":="
			}
			l.lastToken = COLON
			return l.pos, l.lastToken, ":"

		case '/':
			if l.match('/') {
				str := l.readLine()
				pos := l.pos
				l.resetPosition() // B/c we read the whole line we need to reset pos

				l.lastToken = LINECOMMENT
				return pos, LINECOMMENT, str
			}
			l.lastToken = DIV
			return l.pos, DIV, "/"

			//--------------------------------------------------------------------------------
			// - Multi character
			//--------------------------------------------------------------------------------
		case '"':
			startPos := l.pos

			lit := l.lexString('"')
			l.lastToken = STRING
			return startPos, STRING, lit
		case '\'':
			startPos := l.pos

			lit := l.lexString('\'')
			l.lastToken = CHARLIT
			return startPos, CHARLIT, lit
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				lit, tokType := l.lexNumber()
				l.lastToken = tokType
				return startPos, l.lastToken, lit
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				l.lastToken = IDENT

				kw, ok := keywordList[lit]
				if ok {
					l.lastToken = kw
				}

				return startPos, l.lastToken, lit
			} else {
				l.lastToken = ILLEGAL
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

// lexInt scans the input until the end of an integer and then returns the
// literal.
func (l *Lexer) lexNumber() (string, TokenType) {
	tokType := INT
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit, tokType
			}
		}

		if r == '.' {
			tokType = FLOAT
		}

		l.pos.column++
		ok := unicode.IsDigit(r) || r == '.' || r == '_'
		if ok {
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return lit, tokType
		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}
		l.pos.column++

		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexString(border rune) string {
	var lit = string(border)
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}
		l.pos.column++

		if r == '\n' {
			// goto next line
			l.resetPosition()
		}

		if r == border {
			// end string
			lit = lit + string(r)
			return lit
		}

		lit = lit + string(r)
	}
}

func (l *Lexer) readLine() string {
	str, err := l.reader.ReadString('\n') // TODO: If you do this with readrune, then you can use backup
	if err == io.EOF {
		return str
	}

	return str[:len(str)-1]
}

// func (l *Lexer) peek() rune {
// 	r, _ err := l.reader.ReadRune()
// 	if err != nil { panic(err) }

// 	err := l.reader.UnreadRune()
// 	if err != nil { panic(err) }
// }

func (l *Lexer) match(nextRune rune) bool {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return false
		}
	}

	matches := r == nextRune
	if !matches {
		// Unread it if we couldn't find a match
		err := l.reader.UnreadRune()
		if err != nil { panic(err) }
	}

	l.pos.column++
	return matches
}
