pub const Engine = @import("engine.zig");

pub const Scheme = struct {
    args: []const Value,
};

pub const Value = union(enum) {
    int: i32,
    char: u8,
    float: f32,
    scheme: Scheme,
    boolean: bool,
    string: []const u8,
    symbol: []const u8,
    vector: []Value,
    pair: []Value,
};
