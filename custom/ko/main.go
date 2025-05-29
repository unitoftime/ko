package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const BuildDirectory = "./out/"

// or GOOS GOARCH
// <arch>-<sub>-<os>-<abi>
// var targets = map[string]string{
// 	"native":
// }

const (
	CmdBuild = "build"
	CmdRun = "run"
)

func getArg(idx int) string {
	if len(os.Args) < (idx+1) {
		printHelp()
		return ""
	}
	return os.Args[idx]
}

func printHelp() {
	fmt.Println("Usage: ko <cmd>")
	fmt.Println("Commands:")
	fmt.Println("- build: just build the application")
	fmt.Println("- run: builds and then runs the application")
	os.Exit(1)
}

func main() {
	cmd := os.Args[1]

	switch cmd {
	case CmdBuild:
		compile(getArg(2))
	case CmdRun:
		compile(getArg(2))
		localCmd(BuildDirectory, "./run.bin")
	}

}

func localCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}

func compile(inputFile string) {
	file, err := os.Open(inputFile)
	if err != nil { panic(err) }
	defer file.Close()

	tokens := make([]Token, 0)
	lexer := NewLexer(inputFile, file)
	for {
		pos, tok, lit := lexer.Lex()

		tokens = append(tokens, Token{pos, tok, lit})
		// fmt.Printf("./%s:%d:%d\t%s\t%s\n", "test.txt", pos.line, pos.column, tok, lit)
		if tok == EOF {
			break
		}
	}

	tokenList := &Tokens{list: tokens}
	parser := NewParser(tokenList)
	result := parser.Parse(inputFile)

	resolver := NewResolver()
	resolver.Resolve(result)

	buf := &genBuf{
		buf: new(bytes.Buffer),
	}
	buf.Generate(result)

	err = os.MkdirAll(BuildDirectory, 0700)
	if err != nil { panic(err) }

	// fmt.Println(buf.String())
	err = os.WriteFile(BuildDirectory + "main.c", buf.buf.Bytes(), 0644)
	if err != nil { panic(err) }

	// FLAG_BASIC=-Wall -Wextra -Werror
	// FLAG_UB=-Wpedantic -fsanitize=undefined -fsanitize=address -fno-omit-frame-pointer
	// FLAG_STRICT=-Wshadow -Wstrict-prototypes -Wpointer-arith -Wcast-align \
	// -Wwrite-strings -Wswitch-enum -Wunreachable-code \
	// -Wmissing-prototypes -Wdouble-promotion -Wformat=2

	// FLAGS=-g ${FLAG_BASIC} ${FLAG_UB} ${FLAG_STRICT}

	cc := "gcc"
	args := []string{
		"./out/main.c",
		"-g",
		"-std=c11",

		// Flags
		"-Wall", "-Wextra", "-Werror",
		"-Wpedantic", "-fsanitize=undefined", "-fsanitize=address", "-fno-omit-frame-pointer",
		"-Wshadow", "-Wstrict-prototypes", "-Wpointer-arith", "-Wcast-align",
		"-Wwrite-strings", "-Wswitch-enum", "-Wunreachable-code",
		"-Wmissing-prototypes", "-Wdouble-promotion", "-Wformat=2",

		"-o", "./out/run.bin",
	}

	// Build command

	err = localCmd("./", cc, args...)
	if err != nil { panic(err) }
}
