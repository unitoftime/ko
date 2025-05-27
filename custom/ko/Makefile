build:
	go run .

FLAG_BASIC=-Wall -Wextra -Werror
FLAG_UB=-Wpedantic -fsanitize=undefined -fsanitize=address -fno-omit-frame-pointer
FLAG_STRICT=-Wshadow -Wstrict-prototypes -Wpointer-arith -Wcast-align \
-Wwrite-strings -Wswitch-enum -Wunreachable-code \
-Wmissing-prototypes -Wdouble-promotion -Wformat=2

FLAGS=-g ${FLAG_BASIC} ${FLAG_UB} ${FLAG_STRICT}


gcc: build
	cd ./out && gcc ${FLAGS} ./main.c -o run.bin && ./run.bin

tcc: build
	cd ./out && tcc ./main.c -o ./run.bin && ./run.bin

# echo 'int main() { return 42; }' | tcc -o myprog -
# echo 'int main() { return 0; }' | gcc -x c -o myprog -
