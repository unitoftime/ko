const std = @import("std");

pub fn Print(comptime fmt: []const u8, args: anytype) void {
    std.debug.print(fmt, args);
}

pub fn Println(comptime fmt: []const u8, args: anytype) void {
    std.debug.print(fmt, args);
    std.debug.print("\n", .{});
}
