package main

foreign {
	func malloc(size usize) *u8
	func free(ptr *u8)
	// func sizeof[T](T) usize // TODO: bad return type
}

func main() {
	ko_printf("Starting Interpreter\n")
	// chunk := Chunk{[]u8{}, []f64{}, []int{}}
	chunk := Chunk{}

	c := addConstant(&chunk, 1.2)
	writeChunk(&chunk, OpConstant, 123)
	writeChunk(&chunk, c, 123)


	writeChunk(&chunk, OpReturn, 123)
	// TODO: Free chunk

	disassembleChunk(&chunk, "test")
}

// type Opcode enum
const OpReturn = 0
const OpConstant = 1

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
	for offset := 0; offset < len[u8](chunk.code); offset++ {
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
