#ifndef clox_header_h
#define clox_header_h

#define UINT8_COUNT (UINT8_MAX + 1)

typedef struct Obj Obj;
typedef struct ObjString ObjString;

typedef enum {
  OP_CONSTANT,
  OP_NIL,
  OP_TRUE,
  OP_FALSE,
  OP_EQUAL,
  OP_GREATER,
  OP_LESS,
  OP_ADD,
  OP_SUBTRACT,
  OP_MULTIPLY,
  OP_DIVIDE,
  OP_NOT,
  OP_NEGATE,
  OP_PRINT,
  OP_RETURN,
  OP_POP,
  OP_DEFINE_GLOBAL,
  OP_GET_GLOBAL,
  OP_SET_GLOBAL,
  OP_GET_LOCAL,
  OP_SET_LOCAL,
  OP_JUMP_IF_FALSE,
  OP_JUMP,
  OP_LOOP,
  OP_CALL,
} OpCode;

typedef enum {
  OBJ_FUNCTION,
  OBJ_NATIVE,
  OBJ_STRING,
} ObjType;

typedef enum {
  // Single-character tokens.
  TOKEN_LEFT_PAREN, TOKEN_RIGHT_PAREN,
  TOKEN_LEFT_BRACE, TOKEN_RIGHT_BRACE,
  TOKEN_COMMA, TOKEN_DOT, TOKEN_MINUS, TOKEN_PLUS,
  TOKEN_SEMICOLON, TOKEN_SLASH, TOKEN_STAR,
  // One or two character tokens.
  TOKEN_BANG, TOKEN_BANG_EQUAL,
  TOKEN_EQUAL, TOKEN_EQUAL_EQUAL,
  TOKEN_GREATER, TOKEN_GREATER_EQUAL,
  TOKEN_LESS, TOKEN_LESS_EQUAL,
  // Literals.
  TOKEN_IDENTIFIER, TOKEN_STRING, TOKEN_NUMBER,
  // Keywords.
  TOKEN_AND, TOKEN_CLASS, TOKEN_ELSE, TOKEN_FALSE,
  TOKEN_FOR, TOKEN_FUN, TOKEN_IF, TOKEN_NIL, TOKEN_OR,
  TOKEN_PRINT, TOKEN_RETURN, TOKEN_SUPER, TOKEN_THIS,
  TOKEN_TRUE, TOKEN_VAR, TOKEN_WHILE,

  TOKEN_ERROR, TOKEN_EOF
} TokenType;

typedef struct {
  const char* start;
  const char* current;
  int line;
} Scanner;

typedef struct {
  TokenType type;
  const char* start;
  int length;
  int line;
} Token;

typedef enum {
  VAL_BOOL,
  VAL_NIL,
  VAL_NUMBER,
  VAL_OBJ,
} ValueType;

typedef struct {
  ValueType type;
  union {
    bool boolean;
    double number;
    Obj* obj;
  } as;
} Value;

typedef struct {
  int capacity;
  int count;
  Value* values;
} ValueArray;

typedef struct {
  int count;
  int capacity;
  uint8_t* code;
  int* lines;
  ValueArray constants;
} Chunk;

struct Obj {
  ObjType type;
  struct Obj* next;
};

struct ObjString {
  Obj obj;
  int length;
  char* chars;
  uint32_t hash;
};

typedef struct {
  ObjString* key;
  Value value;
} Entry;

typedef struct {
  int count;
  int capacity;
  Entry* entries;
} Table;

typedef struct {
  Token name;
  int depth;
} Local;

typedef enum {
  TYPE_FUNCTION,
  TYPE_SCRIPT
} FunctionType;

typedef struct {
  Obj obj;
  int arity;
  Chunk chunk;
  ObjString* name;
} ObjFunction;

typedef Value (*NativeFn)(int argCount, Value* args);

typedef struct {
  Obj obj;
  NativeFn function;
} ObjNative;

typedef struct Compiler {
  struct Compiler* enclosing;
  ObjFunction* function;
  FunctionType type;

  Local locals[UINT8_COUNT];
  int localCount;
  int scopeDepth;
} Compiler;

Compiler* current = NULL;



typedef struct {
  ObjFunction* function;
  uint8_t* ip;
  Value* slots;
} CallFrame;

#define FRAMES_MAX 64
#define STACK_MAX (FRAMES_MAX * UINT8_COUNT)

typedef struct {
  CallFrame frames[FRAMES_MAX];
  int frameCount;

  Value stack[STACK_MAX];
  Value* stackTop;

  Obj* objects;
  Table globals;
  Table strings;
} VM;

typedef enum {
  INTERPRET_OK,
  INTERPRET_COMPILE_ERROR,
  INTERPRET_RUNTIME_ERROR
} InterpretResult;

void initValueArray(ValueArray* array);
void writeValueArray(ValueArray* array, Value value);
void freeValueArray(ValueArray* array);
void printValue(Value value);
bool valuesEqual(Value a, Value b);

void initScanner(const char* source);
Token scanToken();
ObjFunction* compile(const char* source);
ObjString* copyString(const char* chars, int length);
ObjString* takeString(char* chars, int length);
ObjFunction* newFunction();
void printObject(Value value);

#define OBJ_TYPE(value)        (AS_OBJ(value)->type)
#define IS_STRING(value)       isObjType(value, OBJ_STRING)
#define IS_FUNCTION(value)     isObjType(value, OBJ_FUNCTION)
#define IS_NATIVE(value)       isObjType(value, OBJ_NATIVE)

#define AS_STRING(value)       ((ObjString*)AS_OBJ(value))
#define AS_CSTRING(value)      (((ObjString*)AS_OBJ(value))->chars)
#define AS_FUNCTION(value)     ((ObjFunction*)AS_OBJ(value))
#define AS_NATIVE(value) \
    (((ObjNative*)AS_OBJ(value))->function)


void initVM();
void freeVM();
static InterpretResult run();
void push(Value value);
Value pop();
InterpretResult interpret(const char* source);

void initTable(Table* table);
void freeTable(Table* table);
bool tableGet(Table* table, ObjString* key, Value* value);
bool tableSet(Table* table, ObjString* key, Value value);
bool tableDelete(Table* table, ObjString* key);
void tableAddAll(Table* from, Table* to);
ObjString* tableFindString(Table* table, const char* chars,
                           int length, uint32_t hash);
void freeChunk(Chunk* chunk);

ObjNative* newNative(NativeFn function);


VM vm;

#endif
