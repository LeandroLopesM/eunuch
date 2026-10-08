const std = @import("std");
const Io = std.Io;

const eunuch = @import("eunuch");

pub fn main(init: std.process.Init) !void {
    var arena = std.heap.ArenaAllocator.init(init.gpa);
    defer arena.deinit();

    var lisp_engine: eunuch.Engine = .new(.{
        .allocator = arena.allocator(),
        .error_style = .print_and_report,
        .verbose = true,
    });
    defer lisp_engine.deinit();
}
