# TODO
- [x] ++ --
- [x] :=
- [x] += -=
- [x] Pointers
- [ ] bitshifts: >> <<
- [ ] bitwise: & | ~ ^
- [ ] other unaries: +x ^x *x &x
- [ ] Global variables whose expression is a function call (C doesn't allow this you may need to do an init function or something like that?)
- [ ] Arrays
- [ ] slices
- [ ] Maps
- [ ] Switch statements
- [ ] else if statements
- [ ] enums -> I think just generate compile time constants
- [ ] tagged unions -> compile time completion check
- [ ] Error/Result/Optional Types? Generics? Multiple returns? Tupled returns?
- [ ] alloc/free and general memory safety. References? Ownership? Pointers? Nil Checks?
- [ ] defer (scope based)
- [ ] Packages and package imports
- [ ] Closures (how would you enforce memory safety
- [ ] Address sanitizer: -fsanitize=address
- [ ] valgrind testing mode
- [ ] Static analyzers: clang-tidy, cppcheck
- [ ] Guarded (debug mode) allocator that ensures you dont go out of bounds
- [ ] Something like a goroutine

# Things I'm unsure about
- glibc linking vs musl
- jemalloc vs glibc malloc
- stack unwinding <- impossible? maybe works with libunwind?
  - libunwind
- backtrace() <- platform dependent
  - libunwind

# Just some ideas
## arena { ... } scopes
scopes that switch the default allocator to an arena, with the expectation that all pointers downstream will be allocated to that arena, then you can just statically check and enforce that anything escaping the arena gets deep copied or translated up somehow to the main allocator.
```
func myFunc() string {
    dat := os.ReadAll("myFile.txt")

    var fieldIWanted string
    arena {
        ret := &myStruct
        json.Unmarshal(dat, &ret) // Allocates everything to the arena

        // Do stuff with ret
        fieldIWanted = myStruct.field // Copy externally
    }

    return fieldIWanted
}
```

## C { ... } scopes
Let people inject c code directly inline to their program, which would be nice for slowly porting C code:
- downside, if we ever move to custom backend then we require a C compiler
- for now its helpful though because I just generate C code as IR, so I can just copy-paste?
- wouldn't I need this regardless for bindings though? like for header only libs? Or maybe just require people to always link their .a?

## Something like odin #foreign to make bindings against static libs

## Integer size switches in generated code
```
#include <stdint.h>
#if INTPTR_MAX == INT64_MAX
    // 64-bit
#elif INTPTR_MAX == INT32_MAX
    // 32-bit
#endif
```
