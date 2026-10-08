const Stack = @import("util.zig").Stack;
const core = @import("lib.zig");
const std = @import("std");
const Engine = @This();

pub const Builtin = struct {};

stack: Stack(core.Value),
stack_guards: Stack(i32),

currently_executing: ?core.Scheme,

variables: std.StringHashMap(core.Value),
functions: std.StringHashMap(Builtin),

pub const Options = struct {
    verbose: bool = false,
    error_style: enum { ignore, report, print_and_report } = .report,
    allocator: std.mem.Allocator = std.heap.page_allocator,
};

var options: Options = .{};

pub fn new(opt: Options) Engine {
    options = opt;

    return .{
        .stack = .init(opt.allocator),
        .stack_guards = .init(opt.allocator),
        .currently_executing = null,
        .variables = .init(options.allocator),
        .functions = .init(options.allocator),
    };
}

pub fn deinit(self: *Engine) void {
    self.stack.deinit();
    self.stack_guards.deinit();
    self.variables.deinit();
    self.functions.deinit();
}
