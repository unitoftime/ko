package src

import (
	"fmt"

	rl "github.com/raysan5/raylib"
)

// const mystruct = struct
type mystruct struct {
	x i32
}

func main() {
	const x = add(3, 5)
	fmt.Println("Hello World {d}", mystruct{x: x})

	rl.InitWindow(960, 540, "My Window Name");
	rl.SetTargetFPS(144);
	defer rl.CloseWindow();

	for !rl.WindowShouldClose() {
		rl.BeginDrawing();
		rl.ClearBackground(rl.BLACK);
		rl.EndDrawing();
	}

}

// export fn add(a: i32, b: i32) i32 {
//     return a + b;
// }
func add(a, b i32) i32 {
	// return a + b
	fmt.Println("ADD", struct{}{})
	return a + b
}


// import (
// 	"fmt"
// 	"math"
// )

// func add(a, b int) int {
// 	// return a + b
// 	fmt.Println("ADD")
// 	return a + b + int(math.Floor(1.0))
// }

// var x int
// var y = 10
// const z = 500.0

// type something struct {
// 	a int
// 	b float64
// }

// var (
// 	a int = 5
// 	b = 6
// )

// const (
// 	c = float32(7)
// 	d = 8
// )

/*
const std = @import("std");

pub fn main() !void {
    // Prints to stderr (it's a shortcut based on `std.io.getStdErr()`)
    std.debug.print("All your {s} are belong to us.\n", .{"codebase"});

    // stdout is for the actual output of your application, for example if you
    // are implementing gzip, then only the compressed bytes should be sent to
    // stdout, not any debugging messages.
    const stdout_file = std.io.getStdOut().writer();
    var bw = std.io.bufferedWriter(stdout_file);
    const stdout = bw.writer();

    try stdout.print("Run `zig build test` to run the tests.\n", .{});

    try bw.flush(); // don't forget to flush!
}

test "simple test" {
    var list = std.ArrayList(i32).init(std.testing.allocator);
    defer list.deinit(); // try commenting this out and see if zig detects the memory leak!
    try list.append(42);
    try std.testing.expectEqual(@as(i32, 42), list.pop());
}
*/

// const fmt = @import("lib/fmt.zig");

// pub fn main() void {
//     fmt.Println("Hello World: {d}", .{5});
// }
