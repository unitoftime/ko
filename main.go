package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// TODO: Flags
var Debug = false // Toggle this to enable/disable debug output
var GenerateFile = true
var GenerateOnly = false

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
	fmt.Println("Compile:", inputFile)
	file, err := os.Open(inputFile)
	if err != nil { panic(err) }
	defer file.Close()

	now := time.Now()

	fmt.Printf("Lexer: ")
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

	fmt.Println(time.Since(now))
	now = time.Now()
	fmt.Printf("Parsing: ")
	tokenList := &Tokens{list: tokens}
	parser := NewParser(tokenList)
	result := parser.Parse(inputFile)

	fmt.Println(time.Since(now))
	now = time.Now()
	fmt.Printf("Resolving: ")
	resolver := NewResolver()
	resolver.Resolve(result)

	result.genericInstantiations = resolver.genericInstantiations

	// DebugTree(result.file, 0)

	// ---

	fmt.Println(time.Since(now))
	now = time.Now()
	fmt.Printf("Generating: ")
	err = os.MkdirAll(BuildDirectory, 0700)
	if err != nil { panic(err) }

	// If generate only mode
	if GenerateOnly {
		debugFile, err := os.Create(BuildDirectory + "main.c")
		if err != nil {
			panic(err)
		}
		defer debugFile.Close()
		fileWriter := bufio.NewWriter(debugFile)

		buf := &genBuf{
			buf: fileWriter,
		}
		buf.Generate(result)

		err = fileWriter.Flush()
		if err != nil { panic(err) }
		return
	}

	cmd := pipeCompile()
	pipeFile, err := cmd.StdinPipe()
	if err != nil { panic(err) }
	pipeWriter := bufio.NewWriter(pipeFile)

	var fileWriter *bufio.Writer
	var finalWriter io.Writer
	// Debug mode: Also write to main.c
	if GenerateFile {
		debugFile, err := os.Create(BuildDirectory + "main.c")
		if err != nil {
			panic(err)
		}
		defer debugFile.Close()
		fileWriter = bufio.NewWriter(debugFile)


    finalWriter = io.MultiWriter(pipeWriter, fileWriter)
	} else {
    finalWriter = io.MultiWriter(pipeWriter)
	}


	err = cmd.Start()
	if err != nil { panic(err) }

	buf := &genBuf{
		// buf: new(bytes.Buffer),
		buf: finalWriter,
	}
	buf.Generate(result)

	err = pipeWriter.Flush()
	if err != nil { panic(err) }
	err = fileWriter.Flush()
	if err != nil { panic(err) }

	err = pipeFile.Close()
	if err != nil { panic(err) }

	// // fmt.Println(buf.String())
	// err = os.WriteFile(BuildDirectory + "main.c", buf.buf.Bytes(), 0644)
	// if err != nil { panic(err) }

	// FLAG_BASIC=-Wall -Wextra -Werror
	// FLAG_UB=-Wpedantic -fsanitize=undefined -fsanitize=address -fno-omit-frame-pointer
	// FLAG_STRICT=-Wshadow -Wstrict-prototypes -Wpointer-arith -Wcast-align \
	// -Wwrite-strings -Wswitch-enum -Wunreachable-code \
	// -Wmissing-prototypes -Wdouble-promotion -Wformat=2

	// FLAGS=-g ${FLAG_BASIC} ${FLAG_UB} ${FLAG_STRICT}

	fmt.Println(time.Since(now))
	now = time.Now()

	err = cmd.Wait()
	if err != nil { panic(err) }

	{
		// cc := "tcc"

		// opt := "-O0"
		// args := []string{
		// 	"./out/main.c",
		// 	"-g",
		// 	"-std=c11",
		// 	opt,

		// 	// Flags
		// 	"-Wall", "-Wextra", "-Werror",
		// 	"-Wpedantic", "-fsanitize=undefined", "-fsanitize=address", "-fno-omit-frame-pointer",
		// 	"-Wshadow", "-Wstrict-prototypes", "-Wpointer-arith", "-Wcast-align",
		// 	"-Wwrite-strings", "-Wswitch-enum", "-Wunreachable-code",
		// 	"-Wmissing-prototypes", "-Wdouble-promotion", "-Wformat=2",

		// 	// DISABLE UNUSED FUNCTION ERRORS
		// 	"-Wno-unused-function",

		// 	"-o", "./out/run.bin",
		// }

		// // Build command

		// err = localCmd("./", cc, args...)
		// if err != nil { panic(err) }
	}

	fmt.Printf("Running C Compiler: ")
	fmt.Println(time.Since(now))
}

func pipeCompile() *exec.Cmd {
	cc := "gcc"
	opt := "-O0"
	args := []string{
		"-g",
		"-std=c11",
		opt,

		// -ftime-report <- timing information in gcc and clang

		// Flags
		"-Wall", "-Wextra", "-Werror",
		"-Wpedantic", "-fsanitize=undefined", "-fsanitize=address", "-fno-omit-frame-pointer",
		"-Wshadow", "-Wstrict-prototypes", "-Wpointer-arith", "-Wcast-align",
		"-Wwrite-strings", "-Wswitch-enum", "-Wunreachable-code",
		"-Wmissing-prototypes", "-Wdouble-promotion", "-Wformat=2",

		// DISABLE UNUSED FUNCTION ERRORS
		"-Wno-unused-function",

		// "-x", // TCC
		"-x", "c", // GCC
		"-o", "./out/run.bin",
		"-",
	}

	// Build command
	cmd := exec.Command(cc, args...)
	cmd.Dir = "./"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
