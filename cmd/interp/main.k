package main

foreign {
	func exit(status int)

	type FILE struct {}
	var stdin *FILE
	// char *fgets(char *str, int n, FILE *stream);
	func fgets(str *char, n int, stream *FILE) *char
	func fopen(filename *char, mode *char) *FILE

	var SEEK_END int
	func fseek(stream *FILE, offset i64, whence int)
	func ftell(stream *FILE) i64
	func rewind(stream *FILE)
	func fread(ptr *u8, size usize, nitems usize, stream *FILE) usize

	// Upon successful completion, fclose() returns 0. Otherwise, it returns EOF and sets errno to indicate the error.
	func fclose(stream *FILE) int

	func malloc(size usize) *u8
	func free(ptr *u8)
	func ko_args() []string // Builtin: Returns arguments from argc, argv
	func ko_cstr(s string) *char
	func sizeof[T](val T) usize // TODO: bad return type
}

const DEBUG_TRACE_EXECUTION = true
const STACK_MAX = 256

func main() {
	ko_printf("Starting\n")

	initVM()

	args := ko_args()

	numArgs := len[string](args)
	ko_printf("NumArgs: %d\n", numArgs)
	switch numArgs {
	case 1:
		ko_printf("== Repl Mode ==\n")
		repl()
	case 2:
		ko_printf("== Interp Mode: %s ==\n", args[1])
		runFile(args[1])
	default:
		ko_printf("Usage: Requires 0 or 1 argument\n")
		exit(64)
		// TODO: os.Exit(64) SetReturn: 64
	}

	// TODO: Free chunk, VM, etc
}

func repl() {
	lineLen := 1024
	line := make[[]char](lineLen)
	// lineLen := int(sizeof(line[0]))
	for ;;; {
		ko_printf("> ")

		val := fgets(&line[0], lineLen, stdin)
		if val == nil {
			ko_printf("\n")
			break
		}

		slicedLine := line
		for i := 0; i < lineLen; i++ {
			if slicedLine[i] == '\0' {
				slice[[]char](&slicedLine, 0, i)
				break
			}
		}

		interpret(slicedLine)
	}

	// TODO: defer free of line slice
}

func runFile(path string) {
	ko_printf("runFile: %s\n", path)
	source := readFile(path)
	result := interpret(source)
	// free(source)

	if result == InterpretCompileError {
		exit(65)
	}
	if result == InterpretRuntimeError {
		exit(70)
	}
}

func readFile(path string) []char {
	ko_printf("readFile: %s\n", path)
	// Get File Size
  file := fopen(ko_cstr(path), ko_cstr("rb"))
	if file == nil {
		ko_printf("Could not open file %s", path)
		exit(74)
	}

  fseek(file, 0, SEEK_END)
  fileSize := ftell(file)
  rewind(file)

	ko_printf("filesSize: %d\n", fileSize)

	// Read the file to a buffer
	var char_size char
	sizeof_char := sizeof(char_size)

	// buffer := make([]char, fileSize + 1)
	buffer := make[[]char](fileSize + 1)
	ko_printf("buffer: %d\n", len[char](buffer))

  bytesRead := fread(&buffer[0], sizeof_char, fileSize, file)
  if i64(bytesRead) < fileSize {
    ko_printf("Could not read file %s.\n", path);
    exit(74)
  }

	ko_printf("bytesRead.1: %d\n", bytesRead)

	buffer[bytesRead] = '\0';

	fclose(file)
	return buffer
}

// func main() {
// 	ko_printf("Starting Interpreter\n")
// 	// chunk := Chunk{[]u8{}, []f64{}, []int{}}
// 	chunk := Chunk{}

// 	c := addConstant(&chunk, 1.2)
// 	writeChunk(&chunk, OpConstant, 123)
// 	writeChunk(&chunk, c, 123)

// 	c = addConstant(&chunk, 3.4)
// 	writeChunk(&chunk, OpConstant, 123)
// 	writeChunk(&chunk, c, 123)

// 	writeChunk(&chunk, OpAdd, 123)

// 	c = addConstant(&chunk, 5.6)
// 	writeChunk(&chunk, OpConstant, 123)
// 	writeChunk(&chunk, c, 123)

// 	writeChunk(&chunk, OpDivide, 123)
// 	writeChunk(&chunk, OpNegate, 123)

// 	writeChunk(&chunk, OpReturn, 123)

// 	// disassembleChunk(&chunk, "test")

// 	initVM()

// 	ret := interpret(&chunk)
// 	ko_printf("Return Value: %d\n", ret)

// 	// TODO: Free chunk, VM, etc
// }

var vm VirtualMachine
func initVM() {
	// Reset stack
	vm.stackTop = &vm.stack[0]
}

func push(val f64) {
	*vm.stackTop = val
	vm.stackTop++
}

func pop() f64 {
	vm.stackTop--
	return *vm.stackTop
}

type VirtualMachine struct {
	chunk *Chunk
	ip *u8
	stack [256]f64
	stackTop *f64
}

// TODO: Typed Enumeration
const InterpretOk = 0
const InterpretCompileError = 1
const InterpretRuntimeError = 2

func interpret(source []char) int {
	ko_printf("\nINTERPRET: %d\n", len[char](source))
	compile(source)


	return InterpretOk
}

func compile(source []char) {
	initScanner(source)

	line := -1
	for ;;; {
		token := scanToken()
		if token.line != line {
			ko_printf("%d ", token.line)
			line = token.line
		} else {
			ko_printf("    | ")
		}
		ko_printf("%d %s\n", token.Type, token.str)

		if token.Type == TOKEN_EOF {
			break
		}
	}
}

type Scanner struct {
	start *char
	current *char
	line int
}

var scanner Scanner
func initScanner(source []char) {
	scanner.start = &source[0]
	scanner.current = scanner.start
	scanner.line = 1
}

type Token struct {
	Type TokenType
	line int
	str string
}

type TokenType enum {

}

// Old Interpret for hand assembled chunks
// func interpret(chunk *Chunk) int {
// 	vm.chunk = chunk
// 	vm.ip = &vm.chunk.code[0]
// 	return run()
// }

func readByte() u8 {
	b := *vm.ip
	vm.ip++
	return b
}
func readConst() f64 {
	return vm.chunk.values[int(readByte())]
}

// type Opcode enum
const OpReturn = 0
const OpConstant = 1
const OpNegate = 2
const OpAdd = 3
const OpSubtract = 4
const OpMultiply = 5
const OpDivide = 6

func run() int {
	for ;;; {
		if DEBUG_TRACE_EXECUTION {
			ko_printf("STK ")
			for slot := &vm.stack[0]; slot < vm.stackTop; slot++ {
				ko_printf("[ ")
				printValue(*slot)
				ko_printf(" ]")
			}
			ko_printf("\n")
			ko_printf("DBG ")
			disassembleInstruction(vm.chunk, int(vm.ip - &vm.chunk.code[0]))
		}

		instruction := readByte()
		switch instruction {
		case OpConstant:
			ko_printf("CONSTANT")
			// printValue(constant)
			// ko_printf("\n")
			constant := readConst()
			push(constant)

		case OpNegate:
			push(-pop())

		case OpAdd:
			b := pop()
			a := pop()
			push(a + b)

		case OpSubtract:
			b := pop()
			a := pop()
			push(a - b)

		case OpMultiply:
			b := pop()
			a := pop()
			push(a * b)

		case OpDivide:
			b := pop()
			a := pop()
			push(a / b)

		case OpReturn:
			ko_printf("POP: ")
			printValue(pop())
			ko_printf("\n")
			return InterpretOk
		}
	}
}

// TODO: Union for values array
type Chunk struct {
	code []u8
	values []f64
	lines []int
}

func writeChunk(chunk *Chunk, dat u8, line int) {
	append[u8](&chunk.code, dat)
	append[int](&chunk.lines, line)
}

func addConstant(chunk *Chunk, value f64) u8 {
	append[f64](&chunk.values, value)
	ret := len[f64](chunk.values) - 1
	return u8(ret)
}

// name string
func disassembleChunk(chunk *Chunk, name string) {
	ko_printf("== %s ==\n", name)
	for offset := 0; offset < int(len[u8](chunk.code)); nil {
		offset = disassembleInstruction(chunk, offset)
	}
}

func disassembleInstruction(chunk *Chunk, offset int) int {
	// Print Offset
	// ko_printf("%04d\n", offset)
	ko_printf("%d ", offset)

	// Print Line number
	if offset > 0 && chunk.lines[offset] == chunk.lines[offset-1] {
		ko_printf("   | ")
	} else {
		ko_printf("%d ", chunk.lines[offset])
	}

	// Print Instruction
	inst := chunk.code[offset]
	switch inst {
	case OpReturn:
		return simpleInstruction("OP_RETURN", offset)
	case OpConstant:
		return constantInstruction("OP_CONSTANT", chunk, offset)
	case OpNegate:
		return simpleInstruction("OP_NEGATE", offset)
	case OpAdd:
		return simpleInstruction("OP_ADD", offset)
	case OpSubtract:
		return simpleInstruction("OP_SUBTRACT", offset)
	case OpMultiply:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OpDivide:
		return simpleInstruction("OP_DIVIDE", offset)
	default:
		ko_printf("Unknown opCode; %d\n", inst)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	ko_printf("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, chunk *Chunk, offset int) int {
	cIdx := int(chunk.code[offset+1])
	// printf("%-16s %4d '", name, constant);
	ko_printf("%s %d ", name, cIdx)
	printValue(chunk.values[cIdx])
	ko_printf("\n")

	return offset + 2
}

func printValue(value f64) {
	ko_printf("%g", value)
}
