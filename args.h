

// Prototype
void __ko_init_args(int argc, const char* argv[]);
__ko___ko_string_slice ko_args(void);

// Setups the global arguments
__ko___ko_string_slice __global_args;
void __ko_init_args(int argc, const char* argv[]) {
  for (int i = 0; i < argc; i++) {
    __ko___ko_string_slice_append(&__global_args, __ko_string_make(argv[i]));
  }
}

__ko___ko_string_slice ko_args(void) {
  return __global_args;
}
