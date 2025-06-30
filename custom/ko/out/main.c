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
typedef struct Chunk Chunk;
bool __ko_Chunk_equality(Chunk a, Chunk b);
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

void __ko_uint8_t_slice_append(__ko_uint8_t_slice* s, uint8_t value);

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

    printf("append: cap: %d len: %d val: %d\n", s->cap, s->len, value);

    s->a[s->len++] = value;
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
#line 21 "./cmd/interp/main.k"
void writeChunk (Chunk* chunk , uint8_t dat );
struct Chunk {
	__ko_uint8_t_slice code;
};
bool __ko_Chunk_equality(Chunk a, Chunk b){
	return ((__ko___ko_uint8_t_slice_equality(a.code, b.code) == true));
}
#line 15 "./cmd/interp/main.k"
int OpReturn = 0;
// package main
#line 9 "./cmd/interp/main.k"
int main (void) {
	printf("Starting Interpreter");
	Chunk chunk = (Chunk){ {0} };
return __mainRet__;
}
#line 21 "./cmd/interp/main.k"
void writeChunk (Chunk* chunk , uint8_t dat ) {
	__ko_uint8_t_slice_append((&(chunk->code)), dat);
}
