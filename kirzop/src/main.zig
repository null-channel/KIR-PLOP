const std = @import("std");
const zap = @import("zap");
const zqlite = @import("zqlite");

// NOTE: this is a super simplified example, just using a hashmap to map
// from HTTP path to request function.
fn dispatch_routes(r: zap.Request) void {
    // dispatch
    if (r.path) |the_path| {
        if (routes.get(the_path)) |foo| {
            foo(r);
            return;
        }
    }
    // or default: present menu
    r.sendBody(
        \\ <html>
        \\   <body>
        \\     <p><a href="/static">static</a></p>
        \\     <p><a href="/dynamic">dynamic</a></p>
        \\   </body>
        \\ </html>
    ) catch return;
}

fn static_site(r: zap.Request) void {
    r.sendBody("<html><body><h1>Hello from STATIC ZAP!</h1></body></html>") catch return;
}

var dynamic_counter: i32 = 0;
fn dynamic_site(r: zap.Request) void {
    dynamic_counter += 1;
    var buf: [128]u8 = undefined;
    const filled_buf = std.fmt.bufPrintZ(
        &buf,
        "<html><body><h1>Hello # {d} from DYNAMIC ZAP!!!</h1></body></html>",
        .{dynamic_counter},
    ) catch "ERROR";
    r.sendBody(filled_buf) catch return;
}

var conn: zqlite.Conn = undefined;
fn setup_db() !void {
    const flags = zqlite.OpenFlags.Create | zqlite.OpenFlags.EXResCode;
    errdefer std.debug.print("last error: {s}\n", .{conn.lastError()});
    conn = try zqlite.open("./test.sqlite", flags);
    // good idea to pass EXResCode to get extended result codes (more detailed error codes)

    try conn.exec("create table if not exists btrees (id INTEGER PRIMARY KEY,json text)", .{});
    try conn.exec("insert into btrees (json) values (?1), (?2)", .{ "Leto", "Ghanima" });
}

// list of data
var data_list = std.ArrayList([]const u8).init(std.heap.page_allocator);
fn put_data(r: zap.Request) void {
    // save run data in database
    // because we failed at sqlite we are using memory for now
    std.debug.print("We got a message to stabe data!\n", .{});
    std.debug.print("Body: {d}\n", .{data_list.items.len});
    const data = if (r.body) |b| b else {
        std.debug.print("No data received\n", .{});
        return;
    };
    data_list.append(data) catch {
        std.debug.print("Append Failed\n", .{});
        return;
    };
}

fn get_data(r: zap.Request) void {
    if (data_list.items.len > 0) {
        std.debug.print("Length: {d}\n", .{data_list.items.len});
        r.sendJson(data_list.items[data_list.items.len - 1]) catch {
            std.debug.print("Failed to send latest data\n", .{});
            return;
        };
    } else {
        r.sendBody("No data available") catch {
            std.debug.print("Failed to send body\n", .{});
            return;
        };
    }
}

fn on_data(r: zap.Request) void {
    const http_method = zap.methodToEnum(r.method);
    switch (http_method) {
        .GET => get_data(r),
        .POST => put_data(r),
        else => r.setStatus(.not_implemented),
    }
}

fn setup_routes(a: std.mem.Allocator) !void {
    routes = std.StringHashMap(zap.HttpRequestFn).init(a);
    try routes.put("/static", static_site);
    try routes.put("/dynamic", dynamic_site);
    try routes.put("/data", on_data);
}

var routes: std.StringHashMap(zap.HttpRequestFn) = undefined;

pub fn main() !void {
    try setup_db();
    try setup_routes(std.heap.page_allocator);
    var listener = zap.HttpListener.init(.{
        .port = 3000,
        .on_request = dispatch_routes,
        .log = true,
    });
    try listener.listen();

    std.debug.print("Listening on 0.0.0.0:3000\n", .{});

    zap.start(.{
        .threads = 2,
        .workers = 2,
    });
}
