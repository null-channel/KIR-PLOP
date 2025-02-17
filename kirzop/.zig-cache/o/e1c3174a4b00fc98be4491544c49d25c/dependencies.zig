pub const packages = struct {
    pub const @"12200223d76ab6cd32f75bc2e31463b0b429bb5b2b6fa4ce8f68dea494ca1ec3398b" = struct {
        pub const build_root = "/home/marek/.cache/zig/p/12200223d76ab6cd32f75bc2e31463b0b429bb5b2b6fa4ce8f68dea494ca1ec3398b";
        pub const build_zig = @import("12200223d76ab6cd32f75bc2e31463b0b429bb5b2b6fa4ce8f68dea494ca1ec3398b");
        pub const deps: []const struct { []const u8, []const u8 } = &.{
        };
    };
    pub const @"12205c870252c9d4a38397809f5388b13dbc6a4550f50d5f214a355f601e38814a67" = struct {
        pub const build_root = "/home/marek/.cache/zig/p/12205c870252c9d4a38397809f5388b13dbc6a4550f50d5f214a355f601e38814a67";
        pub const build_zig = @import("12205c870252c9d4a38397809f5388b13dbc6a4550f50d5f214a355f601e38814a67");
        pub const deps: []const struct { []const u8, []const u8 } = &.{
        };
    };
};

pub const root_deps: []const struct { []const u8, []const u8 } = &.{
    .{ "zap", "12200223d76ab6cd32f75bc2e31463b0b429bb5b2b6fa4ce8f68dea494ca1ec3398b" },
    .{ "zqlite", "12205c870252c9d4a38397809f5388b13dbc6a4550f50d5f214a355f601e38814a67" },
};
