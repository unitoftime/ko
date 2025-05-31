package main

func main() {
	TestHelloWorld()
	TestVariablesAndArithmetic()
	TestUnaryOperators()
	TestFib()
	TestStructs()
	TestStructsNested()
	TestForLoop()
	TestForLoopSimple()
	TestIfStatement()
	TestGlobalVariable()
	TestGlobalVariableStruct()
	TestCallWithStruct()
	TestScopeNesting()
	TestPointers()

	// runBenchmark()


	// printf("BacktraceTest: %d\n", backtraceTest(5))
	// closure()
}

func TestHelloWorld() {
	printf("Hello World\n")
}

func TestVariablesAndArithmetic() {
	var a = 10
	b := 20
	c := a * b + 5
	Assert(c == 205)

	a++
	Assert(a == 10 + 1)

	a++
	a++
	Assert(a == 10 + 1 + 2)

	a++
	a++
	a++
	Assert(a == 10 + 1 + 2 + 3)
	a++

	d := 5
	d += 3
	Assert(d == 8)
	d += 3
	Assert(d == 11)
	d -= 1
	Assert(d == 10)
	d -= 10
	Assert(d == 0)
}

func TestUnaryOperators() {
	var x = 5
	Assert(x == 5)
	Assert(-x == -5)
	Assert(!false)
	Assert(!!true)
	Assert(!!!false)
	Assert(!(!(!(!true))))
	Assert(!(-(x) == (5)))
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
	Assert(Point{} != p2)
}

type Rect struct {
	Min Point
	Max Point
}
func TestStructsNested() {
	r := Rect{Point{1, 2}, Point{3, 4}}
	r2 := Rect{}
	Assert(r.Min.X == 1)
	Assert(r.Min.Y == 2)
	Assert(r.Max.X == 3)
	Assert(r.Max.Y == 4)
	Assert(Rect{} != r)
	Assert(r2 != r)

	type Rect2 struct {
		R Rect
	}
	rr := Rect2{}
	Assert(rr.R.Min.X != r.Min.X)
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

func TestForLoopSimple() {
	var num = 20
	var count = 0
	for i := 0; i < num; i++ {
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
	Assert(p2 == p3)
	Assert(p1 != p2)
	Assert(p1 == p4)
}

func TestScopeNesting() {
	var a = 5;

	// {
	// 	var a = 10;
	// 	Assert(a == 10)
	// 	{
	// 		var a = 15;
	// 		Assert(a == 15)
	// 	}
	// 	Assert(a == 10)
	// }

	{
		var b = 10;
		Assert(b == 10)
		{
			var c = 15;
			Assert(c == 15)
		}
		Assert(b == 10)
	}

	Assert(a == 5)
}

func TestArrays() {
	var len = 8
	var myArray [8]int

	for i := 0; i < len; i++ {
		Assert(myArray[i] == 0)
		myArray[0] = 99
	}

	for i := 0; i < len; i++ {
		Assert(myArray[i] == 99)
	}

	// Clear and recheck
	myArray = [8]int{}
	for i := 0; i < len; i++ {
		Assert(myArray[i] == 0)
	}
}

func TestPointers() {
	y := 5
	var x *int
	x = &y
	Assert(*x == 5)
	printf("Pointer: %d", *x)
}

// func backtraceTest(n int) int {
// 	if n <= 1 {
// 		printBacktrace()
// 		return n
// 	}
// 	return n + backtraceTest(n - 1)
// }

// // TODO: I think when I do closures I may want to just manually handle them, rather than rely on gcc
// // Note: I think this only works in gcc and isn't standard c
// func closure() {
// 	var i = 0
// 	func count() {
// 		i = i + 1
// 	}

// 	count()
// 	count()
// 	count()
// 	printf("Closure Count: %d", i)
// }

