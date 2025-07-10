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
