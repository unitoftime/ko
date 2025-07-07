#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <stdarg.h>

#define Assert(cond) do { \
    if (!(cond)) \
        printf("Test failed: %s at %s:%d\n", #cond, __FILE__, __LINE__); \
    else \
        printf("Test passed: %s\n", #cond); \
} while (0)


#define ko_byte_malloc(size) ((uint8_t*)malloc(size))


// Note: Warning: This doesn't support utf8
typedef struct {
  const char* data;
  size_t len;
} __ko_string;

__ko_string __ko_string_make(const char* cstr);
__ko_string __ko_string_slice(__ko_string s, size_t low, size_t high);
bool __ko_string_cmp(__ko_string a, __ko_string b);
void __ko_string_free(__ko_string s);
void ko_printf(__ko_string format, ...);

__ko_string __ko_string_make(const char* cstr) {
  size_t len = strlen(cstr);
  char* copy = malloc(len);
  memcpy(copy, cstr, len);
  return (__ko_string){.data = copy, .len = len};
}

__ko_string __ko_string_slice(__ko_string s, size_t low, size_t high) {
  if (low > high || high > s.len) {
    fprintf(stderr, "__ko_string out of range\n");
    exit(1);
  }

  return (__ko_string){ .data = s.data + low, .len = high - low };
}

bool __ko_string_cmp(__ko_string a, __ko_string b) {
  if (a.len != b.len) return false;
  return memcmp(a.data, b.data, a.len);

}

void __ko_string_free(__ko_string s) {
  free((void*)s.data);
}

void ko_printf(__ko_string format, ...) {
  va_list args;
  va_start(args, format); // Initialize va_list with the last fixed argument

  const char* end = format.data + format.len;

  while (format.data < end) {
  /* while (*format != '\0') { */
    if (*(format.data) == '%') {
      format.data++; // Move past '%'
      switch (*(format.data)) {
      case 'c': {
        char c = va_arg(args, int); // char promotes to int in va_arg
        putchar(c);
        break;
      }
      case 's': {
        __ko_string s = va_arg(args, __ko_string);
        char* tmp = malloc(s.len + 1);
        memcpy(tmp, s.data, s.len);
        tmp[s.len] = '\0';
        fputs(tmp, stdout);
        free(tmp);
        break;

        /* char* s = va_arg(args, char*); */
        /* fputs(s, stdout); */
        /* break; */
      }
      case 'd': {
        int d = va_arg(args, int);
        printf("%d", d); // Using standard printf for integer conversion
        break;
      }
        // Add more cases for other format specifiers (%f, %x, etc.)
      default:
        putchar('%'); // Print '%' if unknown specifier
        putchar(*(format.data));
        break;
      }
    } else {
      putchar(*(format.data)); // Print regular characters
    }
    format.data++;
  }

  va_end(args); // Clean up va_list
}

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
typedef struct Chunk Chunk;
bool __ko_Chunk_equality(Chunk a, Chunk b);
typedef struct __ko_double_slice __ko_double_slice;

struct __ko_double_slice {
    double* a;
    size_t len;
    size_t cap;
};

// Protos
__ko_double_slice __ko_double_slice_new(size_t capacity);
__ko_double_slice __ko_double_slice_init(const double* values, size_t count);
void __ko_double_slice_free(__ko_double_slice* s);
double __ko_double_slice_get(__ko_double_slice* s, size_t index);
void __ko_double_slice_set(__ko_double_slice* s, size_t index, double value);
bool __ko___ko_double_slice_equality(__ko_double_slice a, __ko_double_slice b);

void __ko_double_slice_append(__ko_double_slice* s, double value);
size_t __ko_double_slice_len(__ko_double_slice s);

__ko_double_slice __ko_double_slice_new(size_t capacity) {
    __ko_double_slice s;
    s.a = (double*)malloc(capacity * sizeof(double));
    s.len = 0;
    s.cap = capacity;
    return s;
}
__ko_double_slice __ko_double_slice_init(const double* values, size_t count) {
  __ko_double_slice s = __ko_double_slice_new(count);
  memcpy(s.a, values, sizeof(double) * count);
  s.len = count;
  return s;
}

bool __ko___ko_double_slice_equality(__ko_double_slice a, __ko_double_slice b) {
     return (a.a == b.a) && (a.len == b.len) && (a.cap == b.cap);
}

void __ko_double_slice_free(__ko_double_slice* s) {
    if (s->a != NULL) {
        free(s->a);
        s->a = NULL;
    }
    s->len = 0;
    s->cap = 0;
}

void __ko_double_slice_append(__ko_double_slice* s, double value) {
    if (s->len >= s->cap) {
        size_t new_cap = s->cap == 0 ? 4 : s->cap * 2;
        double* new_data = (double*)realloc(s->a, new_cap * sizeof(double));
        if (!new_data) {
            fprintf(stderr, "Out of memory in append()\n");
            exit(1);
        }
        s->a = new_data;
        s->cap = new_cap;
    }

    printf("append: cap: %ld len: %ld val: %d\n", s->cap, s->len, value);

    s->a[s->len++] = value;
}

size_t __ko_double_slice_len(__ko_double_slice s) {
     return s.len;
}

double __ko_double_slice_get(__ko_double_slice* s, size_t index) {
    if (index >= s->len) {
        fprintf(stderr, "Index out of bounds in get()\n");
        exit(1);
    }
    return s->a[index];
}

void __ko_double_slice_set(__ko_double_slice* s, size_t index, double value) {
    if (index >= s->len) {
        fprintf(stderr, "Index out of bounds in set()\n");
        exit(1);
    }
    s->a[index] = value;
}

typedef struct __ko_uint8_t_slice __ko_uint8_t_slice;

struct __ko_uint8_t_slice {
    uint8_t* a;
    size_t len;
    size_t cap;
};

// Protos
__ko_uint8_t_slice __ko_uint8_t_slice_new(size_t capacity);
__ko_uint8_t_slice __ko_uint8_t_slice_init(const uint8_t* values, size_t count);
void __ko_uint8_t_slice_free(__ko_uint8_t_slice* s);
uint8_t __ko_uint8_t_slice_get(__ko_uint8_t_slice* s, size_t index);
void __ko_uint8_t_slice_set(__ko_uint8_t_slice* s, size_t index, uint8_t value);
bool __ko___ko_uint8_t_slice_equality(__ko_uint8_t_slice a, __ko_uint8_t_slice b);

void __ko_uint8_t_slice_append(__ko_uint8_t_slice* s, uint8_t value);
size_t __ko_uint8_t_slice_len(__ko_uint8_t_slice s);

__ko_uint8_t_slice __ko_uint8_t_slice_new(size_t capacity) {
    __ko_uint8_t_slice s;
    s.a = (uint8_t*)malloc(capacity * sizeof(uint8_t));
    s.len = 0;
    s.cap = capacity;
    return s;
}
__ko_uint8_t_slice __ko_uint8_t_slice_init(const uint8_t* values, size_t count) {
  __ko_uint8_t_slice s = __ko_uint8_t_slice_new(count);
  memcpy(s.a, values, sizeof(uint8_t) * count);
  s.len = count;
  return s;
}

bool __ko___ko_uint8_t_slice_equality(__ko_uint8_t_slice a, __ko_uint8_t_slice b) {
     return (a.a == b.a) && (a.len == b.len) && (a.cap == b.cap);
}

void __ko_uint8_t_slice_free(__ko_uint8_t_slice* s) {
    if (s->a != NULL) {
        free(s->a);
        s->a = NULL;
    }
    s->len = 0;
    s->cap = 0;
}

void __ko_uint8_t_slice_append(__ko_uint8_t_slice* s, uint8_t value) {
    if (s->len >= s->cap) {
        size_t new_cap = s->cap == 0 ? 4 : s->cap * 2;
        uint8_t* new_data = (uint8_t*)realloc(s->a, new_cap * sizeof(uint8_t));
        if (!new_data) {
            fprintf(stderr, "Out of memory in append()\n");
            exit(1);
        }
        s->a = new_data;
        s->cap = new_cap;
    }

    printf("append: cap: %ld len: %ld val: %d\n", s->cap, s->len, value);

    s->a[s->len++] = value;
}

size_t __ko_uint8_t_slice_len(__ko_uint8_t_slice s) {
     return s.len;
}

uint8_t __ko_uint8_t_slice_get(__ko_uint8_t_slice* s, size_t index) {
    if (index >= s->len) {
        fprintf(stderr, "Index out of bounds in get()\n");
        exit(1);
    }
    return s->a[index];
}

void __ko_uint8_t_slice_set(__ko_uint8_t_slice* s, size_t index, uint8_t value) {
    if (index >= s->len) {
        fprintf(stderr, "Index out of bounds in set()\n");
        exit(1);
    }
    s->a[index] = value;
}


#line 9 "./cmd/interp/main.k"
int main (void);
#line 34 "./cmd/interp/main.k"
void writeChunk (Chunk* chunk , uint8_t dat );
#line 38 "./cmd/interp/main.k"
int addConstant (Chunk* chunk , double value );
#line 44 "./cmd/interp/main.k"
void disassembleChunk (Chunk* chunk , __ko_string name );
#line 51 "./cmd/interp/main.k"
int disassembleInstruction (Chunk* chunk , int offset );
#line 64 "./cmd/interp/main.k"
int simpleInstruction (__ko_string name , int offset );
struct Chunk {
	__ko_uint8_t_slice code;
	__ko_double_slice values;
};
bool __ko_Chunk_equality(Chunk a, Chunk b){
	return ((__ko___ko_uint8_t_slice_equality(a.code, b.code) == true) && (__ko___ko_double_slice_equality(a.values, b.values) == true));
}
#line 25 "./cmd/interp/main.k"
int OpReturn = 0;
#line 26 "./cmd/interp/main.k"
int OpConstant = 1;
// package main
#line 9 "./cmd/interp/main.k"
int main (void) {
	ko_printf(__ko_string_make("Starting Interpreter\n"));
	Chunk chunk = (Chunk){ {0}, {0} };
	int c = addConstant((&chunk), 1.2);
	writeChunk((&chunk), 1);
	writeChunk((&chunk), c);
	writeChunk((&chunk), 0);
	;
	disassembleChunk((&chunk), __ko_string_make("test"));
return __mainRet__;
}
#line 34 "./cmd/interp/main.k"
void writeChunk (Chunk* chunk , uint8_t dat ) {
	__ko_uint8_t_slice_append((&chunk->code), dat);
}
#line 38 "./cmd/interp/main.k"
int addConstant (Chunk* chunk , double value ) {
	__ko_double_slice_append((&chunk->values), value);
	return ((__ko_double_slice_len(chunk->values) - 1));
}
#line 44 "./cmd/interp/main.k"
void disassembleChunk (Chunk* chunk , __ko_string name ) {
	ko_printf(__ko_string_make("== %s ==\n"), name);
	for (int offset = 0; (offset < __ko_uint8_t_slice_len(chunk->code)); (offset++)) {
		offset = disassembleInstruction(chunk, offset);
	};
}
#line 51 "./cmd/interp/main.k"
int disassembleInstruction (Chunk* chunk , int offset ) {
	;
	ko_printf(__ko_string_make("%d\n"), offset);
	uint8_t inst = chunk->code.a[offset];
	switch (inst) {
	case 0:
		return (simpleInstruction(__ko_string_make("OP_RETURN"), offset));
	break;
	default:
		ko_printf(__ko_string_make("Unknown opCode; %d\n"), inst);
		return ((offset + 1));
	break;
	};
}
#line 64 "./cmd/interp/main.k"
int simpleInstruction (__ko_string name , int offset ) {
	ko_printf(__ko_string_make("%s\n"), name);
	return ((offset + 1));
}
