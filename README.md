# polyp

A transparent TCP proxy with load balancing. Protocol-agnostic — forwards raw bytes in both directions without inspecting content.

## Features

- Transparent bidirectional TCP forwarding
- Multiple load balancing algorithms
- Multiple backend targets

## Build & Run

```bash
go build -o polyp .
./polyp -port 8080 -config config.json
```

Or without building:

```bash
go run . -port 8080 -config config.json
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `-port` | `80` | Port to listen on |
| `-config` | `config.json` | Path to config file |

## Configuration

```json
{
  "targets": [
    { "ip": "127.0.0.1", "port": "3000" },
    { "ip": "127.0.0.1", "port": "3001" }
  ],
  "loadBalancer": "rr",
  "hashTableSize": 1024
}
```

**Fields:**

| Field | Required | Description |
|---|---|---|
| `targets` | yes | List of backend servers |
| `loadBalancer` | no | Algorithm (see below). Defaults to `rr` |
| `hashTableSize` | no | Max entries in the IP hash cache. Only used with `iphash` |

## Load Balancing Algorithms

| Value | Algorithm | Description |
|---|---|---|
| `rr` | Round-robin | Cycles through backends in order |
| `iphash` | IP hash | Same client IP always routes to the same backend (sticky sessions) |
| `lc` | _(not yet implemented)_ | Falls back to round-robin |
| `random` | _(not yet implemented)_ | Falls back to round-robin |

## Running Tests

```bash
go test ./... -race
```
