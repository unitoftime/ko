package main

foreign {
	func malloc(size usize) *u8
	func free(ptr *u8)
	// func sizeof[T](T) usize // TODO: bad return type
}

func main() {
	printf("Starting Interpreter")
	chunk := Chunk{[]u8{}}
}

// type Opcode enum
var OpReturn = 0

type Chunk struct {
	code []u8
}

func writeChunk(chunk *Chunk, dat u8) {
	append[u8](&(chunk.code), dat)
}

// name string
func disassembleChunk(chunk *Chunk) {
	for offset := 0; offset < len(chunk.code); offset++ {
		offset = disassembleInstruction(chunk, offset)
	}
}

func disassembleInstruction(chunk *Chunk, offset int) int {
	printf("%04d", offset)
	inst := chunk.code[offset]
	switch inst {
	case OpReturn:
		return simpleInstruction("OP_RETURN", offset)
	default:
		printf("Unknown opCode; %d\n", inst)
		return offset + 1
	}
}

// func simpleInstruction(name string, offset int) int {
// 	printf("%s\n", name)
// 	return offset + 1
// }
