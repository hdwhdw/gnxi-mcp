# GNXI MCP Server

A Model Context Protocol (MCP) server that provides network device management capabilities using gNOI (gRPC Network Operations Interface).

## Features

- Get device time via gNOI system.Time service
- Supports insecure gRPC connections
- Configurable device address and port

## Usage

### Building
Build all components:
```bash
go build -o gnxi-mcp-server ./cmd
go build -o test-client ./cmd/test-client
go build -o demo-client ./cmd/demo-client
```

### Running the Server
The server runs over stdin/stdout (MCP stdio transport):
```bash
./gnxi-mcp-server
```

### Testing
Test that the server works correctly:
```bash
./test-client
```

Demo with a specific device IP:
```bash
./demo-client 192.168.1.100
```

## MCP Tool: get_device_time

Retrieves the current system time from a network device using gNOI.

**Parameters:**
- `host` (required): Device hostname or IP address
- `port` (required): gRPC port number (default: 9339)
- `insecure` (optional): Use insecure connection (default: true)

**Example:**
```json
{
  "host": "192.168.1.1",
  "port": 9339,
  "insecure": true
}
```

## Expected Output

**Success** (when device is reachable):
```
✅ Successfully got device time:
Device 192.168.1.1:9339 time: 2024-11-24T16:30:45Z (timestamp: 1732467045000000000)
```

**Error** (when device is not reachable):
```
❌ Error calling device:
rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing: dial tcp 192.168.1.1:9339: connection refused"
```
