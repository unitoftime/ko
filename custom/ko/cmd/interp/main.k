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

// func initChunk(chunk *Chunk) {
// 	chunk.count = 0
// 	chunk.cap = 0
// 	chunk.code = nil
// }

// func writeChunk(chunk *Chunk, dat u8) {
// 	if chunk.cap < chunk.count + 1 {
// 		oldCap := chunk.cap
// 		chunk.cap = growCap(oldCap)
// 		chunk.code = growArray[u8](chunk.code, oldCap, chunk.cap)
// 	}

// 	chunk.code[chunk.count] = dat
// 	chunk.count++
// }

// func growCap(cap int) int {
// 	if cap < 8 {
// 		return 8
// 	}
// 	return cap * 2
// }

// func growArray[T](pointer *T, oldCount int, newCount int) *T {
// 	return reallocate(pointer, sizeof(T) * oldCount, sizeof(T) * newCount)
// }

// func reallocate[T](ptr *T, old usize, new usize) *T {
// 	return ptr
// }
