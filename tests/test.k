package main

	// builtin.AddIdent("printf", &BuiltinNode{getType(&BasicType{"void", false})})
	// builtin.AddIdent("Assert", &BuiltinNode{getType(&BasicType{"void", false})})
	// builtin.AddIdent("ko_byte_malloc", &BuiltinNode{getType(&BasicType{"uint8_t*", false})})
	// builtin.AddIdent("sizeof", &BuiltinNode{getType(&BasicType{"size_t", false})})

foreign {
	func malloc(size usize) *u8
	func free(ptr *u8)
	// func sizeof[T](T) usize // TODO: bad return type
	// func append[T](slice []T, val T) []T
	// func append(slice []int, val int) []int
}

func main() {
	// GenerateCode()

	TestHelloWorld()
	TestVariablesAndArithmetic()
	TestUnaryOperators()
	TestFib()
	TestStructs()
	TestStructsNested()
	TestForLoop()
	TestForLoopSimple()
	TestIfStatement()
	TestSwitchStatement()
	TestGlobalVariable()
	TestGlobalVariableStruct()
	TestCallWithStruct()
	TestScopeNesting()
	TestPointers()
	TestMalloc()
	TestArrays()
	TestGeneric()
	TestSlice()
	TestString()
	TestConstants()
	TestCasting()
	// runBenchmark()

	// printf("BacktraceTest: %d\n", backtraceTest(5))
	// closure()
}

func GenerateCode() {
	max := 100000
	for i := 0; i < max; i++ {
		ko_printf("type TestStruct%d struct { a int }\n", i)
		ko_printf("func TestFunc%d(val int) int { return val }\n", i)
	}
}

func TestHelloWorld() {
	ko_printf("Hello World\n")
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
		// ko_printf("Fib %d: %d\n", i, fib(i))
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

func TestSwitchStatement() {
	x := 2
	switch x {
	case 1:
		Assert(x == 1)
	case 2:
		Assert(x == 2)
	case 3:
		Assert(x == 3)
	default:
		Assert(x == 555)
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
	var length = 8
	var myArray [8]int

	for i := 0; i < length; i++ {
		Assert(myArray[i] == 0)
		myArray[i] = 99
	}

	for i := 0; i < length; i++ {
		Assert(myArray[i] == 99)
	}

	// Clear and recheck
	myArray = [8]int{}
	for i := 0; i < length; i++ {
		Assert(myArray[i] == 0)
	}
}

func TestSlice() {
	var mySlice = []int{1, 2, 3}
	mySlice[2] = 33
	Assert(mySlice[0] == 1)
	Assert(mySlice[1] == 2)
	Assert(mySlice[2] == 33)
	Assert(len[int](mySlice) == 3)

	append[int](&mySlice, 44)
	Assert(mySlice[3] == 44)
	Assert(len[int](mySlice) == 4)


	ko_printf("%d %d %d %d %d\n", mySlice[0], mySlice[1], mySlice[2], mySlice[3], mySlice[4])
}

func TestPointers() {
	y := 5
	var x *int
	x = &y
	Assert(*x == 5)
	ko_printf("Pointer: %d\n", *x)
}

func TestMalloc() {
	var x = malloc(1)
	*x = 5
	free(x)
	// var y = malloc(8 * sizeof(Point))

	// var x = ko_byte_malloc(8)
	// var y = ko_byte_malloc(8 * sizeof(Point))
}

func genericAdd[T](a T, b T) T {
	var c T
	c = 1
	return a + b + c
	// return a + b
}

// type genericStruct[T] struct {
// 	val T
// }

func TestGeneric() {
	// TODO: indexexpr should follow almost the same syntax as call expr just with bracket
	x := genericAdd[int](1, 2)
	y := genericAdd[int](x, 3)
	ko_printf("GenericAdds: %d\n", y)

	// x := genericAdd<int>(1, 2)
	// x := genericAdd<"hello">(1, 2)

	// x := genericAdd(int)(1, 2)
	// x := genericStruct(int){5}
	// TODO: Need to add ability to infer generic types
	// y := genericAdd(x, 4)
	// z := genericAdd(x, y)

	// genericAdd(5.5, 6.6)
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
// 	ko_printf("Closure Count: %d", i)
// }


func TestString() {
	ko_printf("TestString\n")
	var str = "hello"

	str2 := "world"

	ko_printf("%s %s", str, str2)
}

const MYCONST = 5
// const MyConst2 = MYCONST + 6
func TestConstants() {
	const STRING = "asdf"

	ko_printf("Constants: %d %s\n", MYCONST, STRING)
}

func TestCasting() {
	var x = 5
	y := u8(x)
	ko_printf("y: %d", y)
}
