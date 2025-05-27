#include <stdio.h>

int __mainRet__ = 0;
typedef struct structTest structTest;
#line 7 "test.txt"
void structTestPrint ( structTest st  );
#line 16 "test.txt"
void structInFunc (  );
#line 25 "test.txt"
int fib ( int n  );
#line 32 "test.txt"
int add ( int a , int b  );
#line 45 "test.txt"
int add2 ( int b  );
#line 63 "test.txt"
int main (  );
#line 44 "test.txt"
int globA = 5;
// package main
struct structTest {
	int X;
	int Y;
};
#line 7 "test.txt"
void structTestPrint ( structTest st  ) {
	printf( "structTestPrint: {%d, %d}\n", st.X, st.Y);
	int old = st.X;
	st.X = st.Y;
	st.Y = old;
	printf( "flipped:         {%d, %d}\n", st.X, st.Y);
	;
}
#line 16 "test.txt"
void structInFunc (  ) {
	struct structTest {
		int X;
		int Y;
	};
	;
	structTest st = { 2, 3 };
	printf( "structInFunc: {%d, %d}\n", st.X, st.Y);
}
#line 25 "test.txt"
int fib ( int n  ) {
	if ((n <= 1)) {
		return (n);
	};
	return ((fib( (n - 2)) + fib( (n - 1))));
}
#line 32 "test.txt"
int add ( int a , int b  ) {
	int x = (a + ((b * 5) / 3));
	if ((a < b)) {
		x = (a + b);
	} else {
		x = (a - b);
	};
	return (x);
}
#line 45 "test.txt"
int add2 ( int b  ) {
	return ((globA + b));
}
#line 63 "test.txt"
int main (  ) {
	printf( "hello world, %d\n", add( 1, 2));
	for (int i = 0; (i < 20); i = (i + 1)) {
		printf( "Fib %d: %d\n", i, fib( i));
	};
	printf( "GlobAdd(9): %d\n", add2( 4));
	structTest st = { 1, 2 };
	structTestPrint( st);
	structInFunc( );
	;
return __mainRet__;
}
