#include <stdint.h>
#include <stdio.h>

/* #include \"raylib.h\" */

//--------------------------------------------------------------------------------
//--------------------------------------------------------------------------------
//--------------------------------------------------------------------------------


/* #define UNW_LOCAL_ONLY */
/* #include <libunwind.h> */

/* /\* #define MAX_FRAMES 64 *\/ */

/* /\* void printBacktrace2 (void) { *\/ */
/* /\*   void *buffer[MAX_FRAMES]; *\/ */
/* /\*   int num_frames; *\/ */

/* /\*   // Get backtrace addresses *\/ */
/* /\*   num_frames = unw_backtrace(buffer, MAX_FRAMES); *\/ */

/* /\*   unw_word_t ip, sp; *\/ */
/* /\*   char func_name[256]; *\/ */
/* /\*   unw_word_t offset; *\/ */
/* /\*   printf("Backtrace (%d frames):\n", num_frames); *\/ */
/* /\*   for (int i = 0; i < num_frames; i++) { *\/ */
/* /\*     printf("  %p\n", buffer[i]); *\/ */

/* /\*     unw_get_reg(&buffer[i], UNW_REG_IP, &ip); *\/ */
/* /\*     unw_get_reg(&buffer[i], UNW_REG_SP, &sp); *\/ */

/* /\*     if (unw_get_proc_name(&buffer[i], func_name, sizeof(func_name), &offset) == 0) { *\/ */
/* /\*       printf("  %s (+0x%lx) [ip=0x%lx sp=0x%lx]\n", func_name, offset, ip, sp); *\/ */
/* /\*     } else { *\/ */
/* /\*       printf("  <unable to get function name> [ip=0x%lx sp=0x%lx]\n", ip, sp); *\/ */
/* /\*     } *\/ */

/* /\*   } *\/ */

/* /\*   /\\* // Translate addresses into an array of strings *\\/ *\/ */
/* /\*   /\\* char **symbols = unw_backtrace_symbols(buffer, num_frames); *\\/ *\/ */

/* /\*   /\\* if (symbols == NULL) { *\\/ *\/ */
/* /\*   /\\*   perror("backtrace_symbols"); *\\/ *\/ */
/* /\*   /\\*   exit(EXIT_FAILURE); *\\/ *\/ */
/* /\*   /\\* } *\\/ *\/ */

/* /\*   /\\* printf("Backtrace (%d frames):\n", num_frames); *\\/ *\/ */
/* /\*   /\\* for (int i = 0; i < num_frames; i++) { *\\/ *\/ */
/* /\*   /\\*   printf("  %s\n", symbols[i]); *\\/ *\/ */
/* /\*   /\\* } *\\/ *\/ */

/* /\*   /\\* free(symbols); *\\/ *\/ */
/* /\* } *\/ */

/* void printBacktrace (void) { */
/*   unw_cursor_t cursor; unw_context_t uc; */
/*   unw_word_t ip, sp; */

/*   unw_getcontext(&uc); */
/*   unw_init_local(&cursor, &uc); */
/*   while (unw_step(&cursor) > 0) { */
/*     unw_get_reg(&cursor, UNW_REG_IP, &ip); */
/*     unw_get_reg(&cursor, UNW_REG_SP, &sp); */

/*     unw_word_t ip, sp; */
/*     char func_name[256]; */
/*     unw_word_t offset; */
/*     if (unw_get_proc_name(&cursor, func_name, sizeof(func_name), &offset) == 0) { */
/*       printf("  %s (+0x%lx) [ip=0x%lx sp=0x%lx]\n", func_name, offset, ip, sp); */
/*     } else { */
/*       printf("  <unable to get function name> [ip=0x%lx sp=0x%lx]\n", ip, sp); */
/*     } */

/*     /\* printf ("ip = %lx, sp = %lx\n", (long) ip, (long) sp); *\/ */
/*   } */
/* } */

/* void printBacktrace() {
/*     unw_cursor_t cursor; */
/*     unw_context_t context; */

/*     // Initialize cursor to current frame for local unwinding */
/*     unw_getcontext(&context); */
/*     unw_init_local(&cursor, &context); */

/*     printf("Backtrace:\n"); */
/*     while (unw_step(&cursor) > 0) { */
/*         unw_word_t ip, sp; */
/*         char func_name[256]; */
/*         unw_word_t offset; */

/*         unw_get_reg(&cursor, UNW_REG_IP, &ip); */
/*         unw_get_reg(&cursor, UNW_REG_SP, &sp); */

/*         if (unw_get_proc_name(&cursor, func_name, sizeof(func_name), &offset) == 0) { */
/*             printf("  %s (+0x%lx) [ip=0x%lx sp=0x%lx]\n", func_name, offset, ip, sp); */
/*         } else { */
/*             printf("  -- error: unable to get function name -- [ip=0x%lx sp=0x%lx]\n", ip, sp); */
/*         } */
/*     } */
/* } */

/* void bar() { show_backtrace(); } */
/* void foo() { bar(); } */

/* int main() { */
/*     foo(); */
/*     return 0; */
/* } */

/* //-------------------------------------------------------------------------------- */
/* // Panic and Unwinding */
/* //-------------------------------------------------------------------------------- */
/* typedef struct Exception { */
/*     const char *message; */
/*     int code; */
/* } Exception; */

/* typedef struct CatchFrame { */
/*     void *handler_ip;  // instruction pointer to jump to */
/*     void *stack_pointer; */
/*     struct CatchFrame *prev; */
/* } CatchFrame; */

/* __thread CatchFrame *catch_stack = NULL; */

/* #define TRY(handler_label)                          \ */
/*     CatchFrame frame;                               \ */
/*     frame.prev = catch_stack;                       \ */
/*     catch_stack = &frame;                           \ */
/*     if (!setjmp((jmp_buf)frame.env)) {              \ */
/*         goto handler_label##_start;                 \ */
/*     } else goto handler_label##_handler;            \ */
/*     handler_label##_start: */

/* #define CATCH                                        \ */
/*     catch_stack = catch_stack->prev;                \ */
/*     goto handler_label##_end;                       \ */
/*     handler_label##_handler: */

/* #define END_TRY                                      \ */
/*     handler_label##_end: */

/* void throw_exception(Exception *e) { */
/*     unw_cursor_t cursor; */
/*     unw_context_t uc; */
/*     unw_getcontext(&uc); */
/*     unw_init_local(&cursor, &uc); */

/*     while (unw_step(&cursor) > 0) { */
/*         // inspect frames (e.g., via unw_get_proc_name) */
/*         // locate registered handler (e.g., via some table) */
/*         // get saved IP/SP from handler frame and jump */
/*     } */

/*     fprintf(stderr, "Unhandled exception: %s\n", e->message); */
/*     exit(1); */
/* } */

#define Assert(cond) do { \
    if (!(cond)) \
        printf("Test failed: %s at %s:%d\n", #cond, __FILE__, __LINE__); \
    else \
        printf("Test passed: %s\n", #cond); \
} while (0)


//--------------------------------------------------------------------------------

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>

// Define the struct
typedef struct {
    int id;
    float score;
    char name[32];
} Person;

// Field-by-field comparison
int person_eq(const Person *a, const Person *b) {
    return a->id == b->id &&
           a->score == b->score &&
           strcmp(a->name, b->name) == 0;
}

// Memcmp comparison
int person_eq_mem(const Person *a, const Person *b) {
    return memcmp(a, b, sizeof(Person)) == 0;
}

// Timing function
double benchmark(int (*cmp)(const Person *, const Person *), Person *a, Person *b, int iterations) {
    clock_t start = clock();
    volatile int result = 0; // volatile prevents optimization
    for (int i = 0; i < iterations; i++) {
        result ^= cmp(a, b);
    }
    clock_t end = clock();
    return (double)(end - start) / CLOCKS_PER_SEC;
}

int runBenchmark() {
    Person p1 = {123, 98.6f, "Alice"};
    Person p2 = {123, 98.6f, "Alice"};

    const int iterations = 10000000;

    double time_field = benchmark(person_eq, &p1, &p2, iterations);
    double time_mem = benchmark(person_eq_mem, &p1, &p2, iterations);

    printf("Field-by-field time: %.6f seconds\n", time_field);
    printf("memcmp time:         %.6f seconds\n", time_mem);

    return 0;
}

int __mainRet__ = 0;
typedef struct Point Point;
#line 3 "./tests/test.k"
int main ();
#line 27 "./tests/test.k"
void TestHelloWorld ();
#line 31 "./tests/test.k"
void TestVariablesAndArithmetic ();
#line 38 "./tests/test.k"
uint64_t fib (uint64_t n );
#line 45 "./tests/test.k"
void TestFib ();
#line 56 "./tests/test.k"
void TestStructs ();
#line 63 "./tests/test.k"
void TestForLoop ();
#line 74 "./tests/test.k"
void TestIfStatement ();
#line 83 "./tests/test.k"
void TestGlobalVariable ();
#line 89 "./tests/test.k"
void TestGlobalVariableStruct ();
struct Point {
	int X;
	int Y;
};
;
bool __ko_Point_equality(Point a, Point b) {
	return ((a.X == b.X) && (a.Y == b.Y));
}
#line 82 "./tests/test.k"
int globA = 5;
#line 88 "./tests/test.k"
Point globPoint = (Point){ 2, 3 };
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
	;
	;
	;
	;
	;
return __mainRet__;
}
#line 27 "./tests/test.k"
void TestHelloWorld () {
	printf("Hello World");
}
#line 31 "./tests/test.k"
void TestVariablesAndArithmetic () {
	int a = 10;
	int b = 20;
	int c = ((a * b) + 5);
	Assert((c == 205));
}
#line 38 "./tests/test.k"
uint64_t fib (uint64_t n ) {
	if ((n <= 1)) {
		return (n);
	};
	return ((fib((n - 2)) + fib((n - 1))));
}
#line 45 "./tests/test.k"
void TestFib () {
	Assert((fib(1) == 1));
	Assert((fib(15) == 610));
	Assert((fib(20) == 6765));
}
#line 56 "./tests/test.k"
void TestStructs () {
	Point p1 = (Point){ 1, 2 };
	Point p2 = (Point){ p1.Y, p1.X };
	Assert((p1.X == p2.Y));
	Assert((p1.Y == p2.X));
}
#line 63 "./tests/test.k"
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
#line 74 "./tests/test.k"
void TestIfStatement () {
	if ((5 < 10)) {
		Assert((1 == 1));
	} else {
		Assert((1 == 2));
	};
}
#line 83 "./tests/test.k"
void TestGlobalVariable () {
	uint64_t ret = fib(globA);
	Assert((5 == ret));
}
#line 89 "./tests/test.k"
void TestGlobalVariableStruct () {
	Assert((__ko_Point_equality(globPoint, (Point){ 2, 3 }) == true));
}
