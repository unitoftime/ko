const fmt = @import("lib/fmt.zig");
const rl= @cImport({@cInclude("raylib.h"); @cInclude("raymath.h");  @cInclude("rlgl.h"); });

const mystruct = struct{
x: i32
};


pub fn main () void {

const x = add(3, 5)
;
fmt.Println("Hello World {d}", mystructs{.x = x});
rl.InitWindow(960, 540, "My Window Name");
rl.SetTargetFPS(144);
defer rl.CloseWindow();
while (!rl.WindowShouldClose()){

rl.BeginDrawing();
rl.ClearBackground(rl.BLACK);
rl.EndDrawing();
}

}

fn add (a: i32, b: i32) i32 {

fmt.Println("ADD", .{});
return a + b;
}
