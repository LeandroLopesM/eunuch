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
    error_style: enum { report, print_and_report } = .report,
    allocator: std.mem.Allocator = std.heap.page_allocator,
};

pub const options: Options = .{};

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

pub fn executeFile(self: *Engine, file_path: []const u8) !void {
    const io: std.Io.Threaded = .init(options.allocator, .{});
    defer io.deinit();
    errdefer |err| {
        std.log.err("Failed to read file '{s}' ({any})", .{ file_path, err });
    }

    const file = try std.Io.Dir.cwd().openFile(io, file_path, &.{});
    defer file.close(io);

    const stat = try file.stat(io);

    const owned_buffer = try options.allocator.alloc(stat.size);
    defer options.allocator.free(owned_buffer);

    file.reader(io, &.{}).interface.readSliceAll(owned_buffer[0..]);

    return self.executeStr(owned_buffer);
}

pub fn deinit(self: *Engine) void {
    self.stack.deinit();
    self.stack_guards.deinit();
    self.variables.deinit();
    self.functions.deinit();
}
