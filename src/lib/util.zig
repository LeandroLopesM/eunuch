const std = @import("std");

pub fn Stack(comptime T: type) type {
    return struct {
        const Self = @This();

        raw_stack: []T,
        stack_ptr: i32,
        guard: i32,

        allocator: std.mem.Allocator,

        pub fn init(alloc: std.mem.Allocator) Self {
            return .{
                .raw_stack = &.{},
                .stack_ptr = 0,
                .guard = 0,
                .allocator = alloc,
            };
        }

        pub fn push(self: *Self, val: T) !void {
            if (self.stack_ptr + 1 > self.raw_stack.len) {
                self.raw_stack = self.allocator.realloc(self.raw_stack, (self.raw_stack.len + 1) * 2) catch @panic("OOM");

                return error.StackOverflow;
            }

            self.raw_stack[self.stack_ptr] = val;
            self.stack_ptr += 1;
        }

        pub fn pop(self: *Self) !T {
            if (self.stack_ptr - 1 < self.guard) {
                return error.StackUnderflow;
            }

            self.stack_ptr -= 1;
            return self.raw_stack[self.stack_ptr];
        }

        pub fn peek(self: *Self) !T {
            const top = try self.pop();
            self.stack_ptr += 1;

            return top;
        }

        pub fn deinit(self: *Self) void {
            self.allocator.free(self.raw_stack);
        }
    };
}
