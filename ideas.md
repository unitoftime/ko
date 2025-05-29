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
