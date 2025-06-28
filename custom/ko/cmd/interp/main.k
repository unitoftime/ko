package main

foreign {
	func malloc(size usize) *u8
	func free(ptr *u8)
	// func sizeof[T](T) usize // TODO: bad return type
}

func main() {
	printf("Starting Interpreter")
	chunk := []u8{}
}

// type Opcode enum
var OpReturn = 0

// type Chunk struct {
// 	count int
// 	cap int
// 	code *u8
// }

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
