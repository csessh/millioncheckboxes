# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is my attempt at replicating https://onemillioncheckboxes.com/. I thought this very simple and fun idea presents so many different challenges from the technical system design perspective.

Let's build a web based interactive game:
* There are 1 million check-boxes
* Multiple players can connect to this website and simultaneously toggle checkbox.
* When a checkbox is toggled, it is immediately displayed to all connected users.
   * Player A click an empty box, it is checked
   * Player B clocks on the same box soon after, it is then unchecked
   * Its status is displayed simultaneously to every users

We go from there.

## Architecture

- **Entry point**: `cmd/server/main.go` - HTTP server with WebSocket endpoint at `/ws`
- **Redis wrapper**: `internal/redis/redis.go` - Provides connection management and basic operations
- **Communication protocol**: JSON-based messaging over WebSocket with command structure (`cmd`, `key`, `value`)
- **Data pattern**: Handles checkbox-like state (`cb-0` to `cb-99` keys with "true"/"false" values)

## Development Commands

### Running the server
```bash
go run cmd/server/main.go
```

### Building the project
```bash
go build -o server cmd/server/main.go
```

### Managing dependencies
```bash
go mod tidy    # Clean up dependencies
go mod download  # Download dependencies
```

### Testing
```bash
go test ./...  # Run all tests
```

## Dependencies

- **Gorilla WebSocket** (`github.com/gorilla/websocket v1.5.3`): WebSocket implementation
- **Go-Redis** (`github.com/redis/go-redis/v9 v9.12.1`): Redis client library

## Key Implementation Details

- **Redis Connection**: Connects to `localhost:6379` with no authentication by default
- **WebSocket Origin Check**: Currently allows all origins (`CheckOrigin: true`)
- **Data Expiration**: Redis keys are set with 24-hour expiration
- **Message Protocol**:
  - `SET` command updates Redis key-value pairs
  - Server broadcasts existing "true" values for keys `cb-0` through `cb-99` on connection
- **Connection Handling**: Basic WebSocket upgrade with JSON message parsing

## Development Notes

- The server runs on port 8080
- Redis must be running on localhost:6379 before starting the server
- WebSocket endpoint: `ws://localhost:8080/ws`
- Current implementation handles up to 100 checkbox states (`cb-0` to `cb-99`)
