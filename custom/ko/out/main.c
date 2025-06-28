#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

#define Assert(cond) do { \
    if (!(cond)) \
        printf("Test failed: %s at %s:%d\n", #cond, __FILE__, __LINE__); \
    else \
        printf("Test passed: %s\n", #cond); \
} while (0)


#define ko_byte_malloc(size) ((uint8_t*)malloc(size))

/* #include <stdio.h> */
/* #include <string.h> */
/* #include <stdlib.h> */
/* #include <time.h> */

/* // Define the struct */
/* typedef struct { */
/*     int id; */
/*     float score; */
/*     char name[32]; */
/* } Person; */

/* // Field-by-field comparison */
/* int person_eq(const Person *a, const Person *b) { */
/*     return a->id == b->id && */
/*            a->score == b->score && */
/*            strcmp(a->name, b->name) == 0; */
/* } */

/* // Memcmp comparison */
/* int person_eq_mem(const Person *a, const Person *b) { */
/*     return memcmp(a, b, sizeof(Person)) == 0; */
/* } */

/* // Timing function */
/* double benchmark(int (*cmp)(const Person *, const Person *), Person *a, Person *b, int iterations) { */
/*     clock_t start = clock(); */
/*     volatile int result = 0; // volatile prevents optimization */
/*     for (int i = 0; i < iterations; i++) { */
/*         result ^= cmp(a, b); */
/*     } */
/*     clock_t end = clock(); */
/*     return (double)(end - start) / CLOCKS_PER_SEC; */
/* } */

/* int runBenchmark() { */
/*     Person p1 = {123, 98.6f, "Alice"}; */
/*     Person p2 = {123, 98.6f, "Alice"}; */

/*     const int iterations = 10000000; */

/*     double time_field = benchmark(person_eq, &p1, &p2, iterations); */
/*     double time_mem = benchmark(person_eq_mem, &p1, &p2, iterations); */

/*     printf("Field-by-field time: %.6f seconds\n", time_field); */
/*     printf("memcmp time:         %.6f seconds\n", time_mem); */

/*     return 0; */
/* } */

int __mainRet__ = 0;
typedef struct Point Point;
bool __ko_Point_equality(Point a, Point b);
typedef struct Rect Rect;
bool __ko_Rect_equality(Rect a, Rect b);
typedef struct __ko_int_slice __ko_int_slice;

struct __ko_int_slice {
    int* a;
    size_t len;
    size_t cap;
};

// Protos
__ko_int_slice __ko_int_slice_new(size_t capacity);
__ko_int_slice __ko_int_slice_init(const int* values, size_t count);
void __ko_int_slice_free(__ko_int_slice* s);
void __ko_int_slice_append(__ko_int_slice* s, int value);
int __ko_int_slice_get(__ko_int_slice* s, size_t index);
void __ko_int_slice_set(__ko_int_slice* s, size_t index, int value);

__ko_int_slice __ko_int_slice_new(size_t capacity) {
    __ko_int_slice s;
    s.a = (int*)malloc(capacity * sizeof(int));
    s.len = 0;
    s.cap = capacity;
    return s;
}
__ko_int_slice __ko_int_slice_init(const int* values, size_t count) {
  __ko_int_slice s = __ko_int_slice_new(count);
  memcpy(s.a, values, sizeof(int) * count);
  s.len = count;
  return s;
}

void __ko_int_slice_free(__ko_int_slice* s) {
    if (s->a != NULL) {
        free(s->a);
        s->a = NULL;
    }
    s->len = 0;
    s->cap = 0;
}

void __ko_int_slice_append(__ko_int_slice* s, int value) {
    if (s->len >= s->cap) {
        size_t new_cap = s->cap == 0 ? 4 : s->cap * 2;
        int* new_data = (int*)realloc(s->a, new_cap * sizeof(int));
        if (!new_data) {
            fprintf(stderr, "Out of memory in append()\n");
            exit(1);
        }
        s->a = new_data;
        s->cap = new_cap;
    }

    printf("append: cap: %d len: %d val: %d\n", s->cap, s->len, value);

    s->a[s->len++] = value;
}

int __ko_int_slice_get(__ko_int_slice* s, size_t index) {
    if (index >= s->len) {
        fprintf(stderr, "Index out of bounds in get()\n");
        exit(1);
    }
    return s->a[index];
}

void __ko_int_slice_set(__ko_int_slice* s, size_t index, int value) {
    if (index >= s->len) {
        fprintf(stderr, "Index out of bounds in set()\n");
        exit(1);
    }
    s->a[index] = value;
}

typedef struct __ko_8int_arr __ko_8int_arr;
struct __ko_8int_arr {
	int a[8];
};

#line 16 "./tests/test.k"
int main (void);
#line 43 "./tests/test.k"
void GenerateCode (void);
#line 51 "./tests/test.k"
void TestHelloWorld (void);
#line 55 "./tests/test.k"
void TestVariablesAndArithmetic (void);
#line 85 "./tests/test.k"
void TestUnaryOperators (void);
#line 96 "./tests/test.k"
uint64_t fib (uint64_t n );
#line 103 "./tests/test.k"
void TestFib (void);
#line 114 "./tests/test.k"
void TestStructs (void);
#line 126 "./tests/test.k"
void TestStructsNested (void);
#line 143 "./tests/test.k"
void TestForLoop (void);
#line 154 "./tests/test.k"
void TestForLoopSimple (void);
#line 164 "./tests/test.k"
void TestIfStatement (void);
#line 173 "./tests/test.k"
void TestGlobalVariable (void);
#line 179 "./tests/test.k"
void TestGlobalVariableStruct (void);
#line 183 "./tests/test.k"
Point reverse (Point val );
#line 189 "./tests/test.k"
void TestCallWithStruct (void);
#line 200 "./tests/test.k"
void TestScopeNesting (void);
#line 226 "./tests/test.k"
void TestArrays (void);
#line 246 "./tests/test.k"
void TestSlice (void);
#line 258 "./tests/test.k"
void TestPointers (void);
#line 266 "./tests/test.k"
void TestMalloc (void);
#line 287 "./tests/test.k"
void TestGeneric (void);
#line 276 "./tests/test.k"
int func_genericAdd___int_int_int (int a , int b );
struct Point {
	int X;
	int Y;
};
bool __ko_Point_equality(Point a, Point b){
	return ((a.X == b.X) && (a.Y == b.Y));
}
struct Rect {
	Point Min;
	Point Max;
};
bool __ko_Rect_equality(Rect a, Rect b){
	return ((__ko_Point_equality(a.Min, b.Min) == true) && (__ko_Point_equality(a.Max, b.Max) == true));
}
#line 276 "./tests/test.k"
int func_genericAdd___int_int_int (int a , int b ) {
	int c = 0;
	c = 1;
	return (((a + b) + c));
	;
}
;
#line 172 "./tests/test.k"
int globA = 5;
#line 178 "./tests/test.k"
Point globPoint = { 2, 3 };
// package main
#line 16 "./tests/test.k"
int main (void) {
	;
	TestHelloWorld();
	TestVariablesAndArithmetic();
	TestUnaryOperators();
	TestFib();
	TestStructs();
	TestStructsNested();
	TestForLoop();
	TestForLoopSimple();
	TestIfStatement();
	TestGlobalVariable();
	TestGlobalVariableStruct();
	TestCallWithStruct();
	TestScopeNesting();
	TestPointers();
	TestMalloc();
	TestArrays();
	TestGeneric();
	TestSlice();
	;
	;
	;
return __mainRet__;
}
#line 43 "./tests/test.k"
void GenerateCode (void) {
	int max = 100000;
	for (int i = 0; (i < max); (i++)) {
		printf("type TestStruct%d struct { a int }\n", i);
		printf("func TestFunc%d(val int) int { return val }\n", i);
	};
}
#line 51 "./tests/test.k"
void TestHelloWorld (void) {
	printf("Hello World\n");
}
#line 55 "./tests/test.k"
void TestVariablesAndArithmetic (void) {
	int a = 10;
	int b = 20;
	int c = ((a * b) + 5);
	Assert((c == 205));
	(a++);
	Assert((a == (10 + 1)));
	(a++);
	(a++);
	Assert((a == ((10 + 1) + 2)));
	(a++);
	(a++);
	(a++);
	Assert((a == (((10 + 1) + 2) + 3)));
	(a++);
	int d = 5;
	d += 3;
	Assert((d == 8));
	d += 3;
	Assert((d == 11));
	d -= 1;
	Assert((d == 10));
	d -= 10;
	Assert((d == 0));
}
#line 85 "./tests/test.k"
void TestUnaryOperators (void) {
	int x = 5;
	Assert((x == 5));
	Assert(((-x) == (-5)));
	Assert((!false));
	Assert((!(!true)));
	Assert((!(!(!false))));
	Assert((!((!((!((!true))))))));
	Assert((!(((-(x)) == (5)))));
}
#line 96 "./tests/test.k"
uint64_t fib (uint64_t n ) {
	if ((n <= 1)) {
		return (n);
	};
	return ((fib((n - 2)) + fib((n - 1))));
}
#line 103 "./tests/test.k"
void TestFib (void) {
	Assert((fib(1) == 1));
	Assert((fib(15) == 610));
	Assert((fib(20) == 6765));
}
#line 114 "./tests/test.k"
void TestStructs (void) {
	Point p1 = (Point){ 1, 2 };
	Point p2 = (Point){ p1.Y, p1.X };
	Assert((p1.X == p2.Y));
	Assert((p1.Y == p2.X));
	Assert((__ko_Point_equality((Point){ 0, 0 }, p2) != true));
}
#line 126 "./tests/test.k"
void TestStructsNested (void) {
	Rect r = (Rect){ (Point){ 1, 2 }, (Point){ 3, 4 } };
	Rect r2 = (Rect){ (Point){ 0, 0 }, (Point){ 0, 0 } };
	Assert((r.Min.X == 1));
	Assert((r.Min.Y == 2));
	Assert((r.Max.X == 3));
	Assert((r.Max.Y == 4));
	Assert((__ko_Rect_equality((Rect){ (Point){ 0, 0 }, (Point){ 0, 0 } }, r) != true));
	Assert((__ko_Rect_equality(r2, r) != true));
	typedef struct Rect2 {
		Rect R;
	} Rect2;
	;
	Rect2 rr = (Rect2){ (Rect){ (Point){ 0, 0 }, (Point){ 0, 0 } } };
	Assert((rr.R.Min.X != r.Min.X));
}
#line 143 "./tests/test.k"
void TestForLoop (void) {
	int num = 20;
	int count = 0;
	for (int i = 0; (i < num); i = (i + 1)) {
		;
		Assert((i == count));
		count = (count + 1);
	};
	Assert((count == num));
}
#line 154 "./tests/test.k"
void TestForLoopSimple (void) {
	int num = 20;
	int count = 0;
	for (int i = 0; (i < num); (i++)) {
		Assert((i == count));
		count = (count + 1);
	};
	Assert((count == num));
}
#line 164 "./tests/test.k"
void TestIfStatement (void) {
	if ((5 < 10)) {
		Assert((1 == 1));
	} else {
		Assert((1 == 2));
	};
}
#line 173 "./tests/test.k"
void TestGlobalVariable (void) {
	uint64_t ret = fib(globA);
	Assert((5 == ret));
}
#line 179 "./tests/test.k"
void TestGlobalVariableStruct (void) {
	Assert((__ko_Point_equality(globPoint, (Point){ 2, 3 }) == true));
}
#line 183 "./tests/test.k"
Point reverse (Point val ) {
	int tmp = val.X;
	val.X = val.Y;
	val.Y = tmp;
	return (val);
}
#line 189 "./tests/test.k"
void TestCallWithStruct (void) {
	Point p1 = (Point){ 1, 2 };
	Point p2 = (Point){ 2, 1 };
	Point p3 = reverse(p1);
	Point p4 = reverse(p3);
	Assert((__ko_Point_equality(p2, p3) == true));
	Assert((__ko_Point_equality(p1, p2) != true));
	Assert((__ko_Point_equality(p1, p4) == true));
}
#line 200 "./tests/test.k"
void TestScopeNesting (void) {
	int a = 5;
	;
	;
	;
	;
	;
	;
	;
	;
	;
	{
		int b = 10;
		Assert((b == 10));
		{
			int c = 15;
			Assert((c == 15));
		}
		;
		Assert((b == 10));
	}
	;
	Assert((a == 5));
}
#line 226 "./tests/test.k"
void TestArrays (void) {
	int len = 8;
	__ko_8int_arr myArray = {0};
	for (int i = 0; (i < len); (i++)) {
		Assert((myArray.a[i] == 0));
		myArray.a[i] = 99;
	};
	for (int i = 0; (i < len); (i++)) {
		Assert((myArray.a[i] == 99));
	};
	;
	myArray = (__ko_8int_arr){0};
	for (int i = 0; (i < len); (i++)) {
		Assert((myArray.a[i] == 0));
	};
}
#line 246 "./tests/test.k"
void TestSlice (void) {
	__ko_int_slice mySlice = __ko_int_slice_init((int[]){1, 2, 3}, 3);
	mySlice.a[2] = 33;
	Assert((mySlice.a[0] == 1));
	Assert((mySlice.a[1] == 2));
	Assert((mySlice.a[2] == 33));
	__ko_int_slice_append((&mySlice), 44);
	Assert((mySlice.a[3] == 44));
	printf("%d %d %d %d %d", mySlice.a[0], mySlice.a[1], mySlice.a[2], mySlice.a[3], mySlice.a[4]);
}
#line 258 "./tests/test.k"
void TestPointers (void) {
	int y = 5;
	int* x = NULL;
	x = (&y);
	Assert(((*x) == 5));
	printf("Pointer: %d\n", (*x));
}
#line 266 "./tests/test.k"
void TestMalloc (void) {
	uint8_t* x = malloc(1);
	(*x) = 5;
	free(x);
	;
	;
	;
}
#line 287 "./tests/test.k"
void TestGeneric (void) {
	;
	int x = func_genericAdd___int_int_int(1, 2);
	int y = func_genericAdd___int_int_int(x, 3);
	printf("GenericAdds: %d\n", y);
	;
	;
	;
	;
	;
	;
	;
	;
}
