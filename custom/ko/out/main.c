#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

#define Assert(cond) do { \
    if (!(cond)) \
        printf("Test failed: %s at %s:%d\n", #cond, __FILE__, __LINE__); \
    else \
        printf("Test passed: %s\n", #cond); \
} while (0)


//--------------------------------------------------------------------------------

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
#line 3 "./tests/test.k"
int main ();
#line 20 "./tests/test.k"
void TestHelloWorld ();
#line 24 "./tests/test.k"
void TestVariablesAndArithmetic ();
#line 31 "./tests/test.k"
uint64_t fib (uint64_t n );
#line 38 "./tests/test.k"
void TestFib ();
#line 49 "./tests/test.k"
void TestStructs ();
#line 56 "./tests/test.k"
void TestForLoop ();
#line 67 "./tests/test.k"
void TestIfStatement ();
#line 76 "./tests/test.k"
void TestGlobalVariable ();
#line 82 "./tests/test.k"
void TestGlobalVariableStruct ();
#line 86 "./tests/test.k"
Point reverse (Point val );
#line 92 "./tests/test.k"
void TestCallWithStruct ();
struct Point {
	int X;
	int Y;
};
bool __ko_Point_equality(Point a, Point b){
	return ((a.X == b.X) && (a.Y == b.Y));
}
#line 75 "./tests/test.k"
int globA = 5;
#line 81 "./tests/test.k"
Point globPoint = { 2, 3 };
// package main
#line 3 "./tests/test.k"
int main () {
	TestHelloWorld();
	TestVariablesAndArithmetic();
	TestFib();
	TestStructs();
	TestForLoop();
	TestIfStatement();
	TestGlobalVariable();
	TestGlobalVariableStruct();
	;
	;
	;
return __mainRet__;
}
#line 20 "./tests/test.k"
void TestHelloWorld () {
	printf("Hello World");
}
#line 24 "./tests/test.k"
void TestVariablesAndArithmetic () {
	int a = 10;
	int b = 20;
	int c = ((a * b) + 5);
	Assert((c == 205));
}
#line 31 "./tests/test.k"
uint64_t fib (uint64_t n ) {
	if ((n <= 1)) {
		return (n);
	};
	return ((fib((n - 2)) + fib((n - 1))));
}
#line 38 "./tests/test.k"
void TestFib () {
	Assert((fib(1) == 1));
	Assert((fib(15) == 610));
	Assert((fib(20) == 6765));
}
#line 49 "./tests/test.k"
void TestStructs () {
	Point p1 = (Point){ 1, 2 };
	Point p2 = (Point){ p1.Y, p1.X };
	Assert((p1.X == p2.Y));
	Assert((p1.Y == p2.X));
}
#line 56 "./tests/test.k"
void TestForLoop () {
	int num = 20;
	int count = 0;
	for (int i = 0; (i < num); i = (i + 1)) {
		;
		Assert((i == count));
		count = (count + 1);
	};
	Assert((count == num));
}
#line 67 "./tests/test.k"
void TestIfStatement () {
	if ((5 < 10)) {
		Assert((1 == 1));
	} else {
		Assert((1 == 2));
	};
}
#line 76 "./tests/test.k"
void TestGlobalVariable () {
	uint64_t ret = fib(globA);
	Assert((5 == ret));
}
#line 82 "./tests/test.k"
void TestGlobalVariableStruct () {
	Assert((__ko_Point_equality(globPoint, (Point){ 2, 3 }) == true));
}
#line 86 "./tests/test.k"
Point reverse (Point val ) {
	int tmp = val.X;
	val.X = val.Y;
	val.Y = tmp;
	return (val);
}
#line 92 "./tests/test.k"
void TestCallWithStruct () {
	Point p1 = (Point){ 1, 2 };
	Point p2 = (Point){ 2, 1 };
	Point p3 = reverse(p1);
	Point p4 = reverse(p3);
	Assert((__ko_Point_equality(p1, p3) == true));
	Assert((__ko_Point_equality(p1, p2) != true));
	Assert((__ko_Point_equality(p2, p4) == true));
}
