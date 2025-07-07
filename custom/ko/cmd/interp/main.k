package main

foreign {
	func malloc(size usize) *u8
	func free(ptr *u8)
	// func sizeof[T](T) usize // TODO: bad return type
}

func main() {
	ko_printf("Starting Interpreter\n")
	chunk := Chunk{[]u8{}, []f64{}}

	c := addConstant(&chunk, 1.2)
	writeChunk(&chunk, OpConstant)
	writeChunk(&chunk, c)


	writeChunk(&chunk, OpReturn)
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
}

func writeChunk(chunk *Chunk, dat u8) {
	append[u8](&chunk.code, dat)
}

func addConstant(chunk *Chunk, value f64) int {
	append[f64](&chunk.values, value)
	return len[f64](chunk.values) - 1
}

// name string
func disassembleChunk(chunk *Chunk, name string) {
	ko_printf("== %s ==\n", name)
	for offset := 0; offset < len[u8](chunk.code); offset++ {
		offset = disassembleInstruction(chunk, offset)
	}
}

func disassembleInstruction(chunk *Chunk, offset int) int {
	// ko_printf("%04d\n", offset)
	ko_printf("%d\n", offset)
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
	ko_printf("%s %d", name, cIdx)
	printValue(chunk.values[cIdx])
	ko_printf("\n")

	return offset + 2
}

func printValue(value f64) {
	printf("%g", value)
}
