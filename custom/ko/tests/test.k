package main

func main() {
	TestHelloWorld()
	TestVariablesAndArithmetic()
	TestFib()
	TestStructs()
	TestForLoop()
	TestIfStatement()
	TestGlobalVariable()
	TestGlobalVariableStruct()

	// runBenchmark()


	// printf("BacktraceTest: %d\n", backtraceTest(5))
	// closure()
}

func TestHelloWorld() {
	printf("Hello World")
}

func TestVariablesAndArithmetic() {
	var a = 10
	var b = 20
	var c = a * b + 5
	Assert(c == 205)
}

func fib(n u64) u64 {
	if n <= 1 {
		return n
	}
	return fib(n - 2) + fib(n - 1)
}

func TestFib() {
	Assert(fib(1) == 1)
	Assert(fib(15) == 610)
	Assert(fib(20) == 6765)
}

type Point struct {
	X int
	Y int
}

func TestStructs() {
	var p1 = Point{1, 2}
	var p2 = Point{p1.Y, p1.X}
	Assert(p1.X == p2.Y)
	Assert(p1.Y == p2.X)
}

func TestForLoop() {
	var num = 20
	var count = 0
	for var i = 0; i < num; i = i + 1 {
		// printf("Fib %d: %d\n", i, fib(i))
		Assert(i == count)
		count = count + 1
	}
	Assert(count == num)
}

func TestIfStatement() {
	if 5 < 10 {
		Assert(1 == 1)
	} else {
		Assert(1 == 2)
	}
}

var globA = 5
func TestGlobalVariable() {
	var ret = fib(globA)
	Assert(5 == ret)
}

var globPoint = Point{2, 3}
func TestGlobalVariableStruct() {
	Assert(globPoint == Point{2, 3})
}

func reverse(val Point) Point {
	var tmp = val.X
	val.X = val.Y
	val.Y = tmp
	return val
}
func TestCallWithStruct() {
	var p1 = Point{1, 2}
	var p2 = Point{2, 1}

	var p3 = reverse(p1)
	var p4 = reverse(p3)
	Assert(p1 == p3)
	Assert(p1 != p2)
	Assert(p2 == p4)
}

// func backtraceTest(n int) int {
// 	if n <= 1 {
// 		printBacktrace()
// 		return n
// 	}
// 	return n + backtraceTest(n - 1)
// }

// // TODO: I think when I do closures I may want to just manually handle them, rather than rely on gcc
// // // Note: I think this only works in gcc and isn't standard c
// // func closure() {
// // 	var i = 0
// // 	func count() {
// // 		i = i + 1
// // 	}

// // 	count()
// // 	count()
// // 	count()
// // 	printf("Closure Count: %d", i)
// // }

