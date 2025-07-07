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
      case 'g': {
        double g = va_arg(args, double);
        printf("%g", g); // Using standard printf for integer conversion
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
