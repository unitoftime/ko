# Interesting Thoughts
- [ ] A token list is kind of like a byte code that we interpret to generate code and typecheck
- [ ] comptime is kind of like a way to add and use variables in our code generation bytecode interpreter
- [ ] Closures:
  - [ ] If they don't escape to the heap, then we can rely on the stack allocated variables (i'm pretty sure)
  - [ ] If they escape to the heap, we could maybe escape and refcount every variable that the closure closes over. I think this would only leak memory if you ended up creating a cycle of closure variables, which is a weird edgecase to think about)
- [ ] Defaults are super nice for defining lots of configuration structs. Like "I want this baseline config, but I want to modify just these fields"

# TODO
- [x] ++ --
- [x] :=
- [x] += -=
- [x] Structs
- [x] Pointers
- [x] Arrays
- [x] Foreign functions <- So I can typewrap stdlib things like malloc,free, etc
- [x] Comptime templating instead?
- [x] Switch statements
- [x] Constant folding
- [ ] Casting types
- [ ] Error/Result/Optional Types? Generics? Multiple returns? Tupled returns?
- [ ] https://nickav.co/posts/0003_wasm_from_scratch
- [ ] variadics
- [ ] Super heavy runtime hotreloading emphasis
- [ ] Global variables whose expression is a function call (C doesn't allow this you may need to do an init function or something like that?)
- [ ] slices <- maybe implement in lang?
- [ ] Maps <- maybe implement in lang?
- [ ] Methods
- [ ] bitshifts: >> << (maybe only allow unsigned types)
- [ ] bitwise: & | ~ ^
- [ ] other unaries: +x ^x
- [ ] else if statements
- [ ] Recursive generic templating - typechecking currntly doesn't work
- [ ] enums -> I think just generate compile time constants. Maybe iota?
- [ ] tagged unions -> compile time completion check
- [ ] alloc/free and general memory safety. References? Ownership? Pointers? Nil Checks?
- [ ] defer (scope based)
- [ ] Packages and package imports
- [ ] Closures (how would you enforce memory safety
- [ ] Address sanitizer: -fsanitize=address
- [ ] valgrind testing mode
- [ ] Static analyzers: clang-tidy, cppcheck
- [ ] Guarded (debug mode) allocator that ensures you dont go out of bounds
- [ ] Something like a goroutine
- [ ] prevent signed integer overflow (on adds and subtracts) - There's builtins in gcc
- [ ] Default functions for structs: Just to have one single place that I can call to return a default struct <- maybe not needed, bc I think {0} works for all my cases
- [ ] When assigning complits and default values, there is probably some places where I can do things like {{0}} or {0},
- [ ] is there a warning flag I can enable to catch implicit casting? So they never happen

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
