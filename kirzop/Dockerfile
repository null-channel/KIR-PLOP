FROM alpine:3.21 AS builder

# Install zig and build dependencies
WORKDIR /app
RUN apk add --no-cache git build-base zig
COPY . .

# Build the application
RUN zig build

FROM alpine:3.19

WORKDIR /app
# Copy the built binary from builder stage
COPY --from=builder /app/zig-out/bin/kirzop .
# Copy the SQLite database file if it exists
COPY --from=builder /app/test.sqlite ./test.sqlite

EXPOSE 3000

CMD ["./kirzop"] 
