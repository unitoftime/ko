const fmt = @import("lib/fmt.zig");
const mystruct = struct{
x: i32
};

pub fn main () void {

const x = add(3, 5);
fmt.Println("Hello World {d}", mystruct{.x = x});
}

fn add (a: i32, b: i32) i32 {

fmt.Println("ADD", .{});
return a + b;
}
