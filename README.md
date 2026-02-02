# Thesis Prototype

Small stateless HTTP microservice on Go: health check and a demo process endpoint. Clean layering (handler → service), config from env, graceful shutdown.

## Requirements

- Go 1.22+

## Project structure

```
thesisPrototype/
├── go.mod
├── README.md
├── cmd/server/
│   └── main.go          # entry point, wiring, server start, shutdown
└── internal/
    ├── config/          # config from environment variables
    ├── logging/         # structured logging (slog)
    ├── handler/         # HTTP handlers (no business logic)
    └── service/         # business logic (process payload)
```

## Build and run

```bash
# build
go build -o server ./cmd/server

# run (default :8080)
./server
```

Or run without building:

```bash
go run ./cmd/server
```

With env:

```bash
HTTP_ADDR=:9090 LOG_LEVEL=debug go run ./cmd/server
```

## Configuration (environment variables)

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | Listen address |
| `HTTP_READ_TIMEOUT` | `10s` | Read timeout |
| `HTTP_WRITE_TIMEOUT` | `10s` | Write timeout |
| `SHUTDOWN_TIMEOUT` | `15s` | Graceful shutdown wait |
| `REQUEST_TIMEOUT` | `30s` | Per-request context timeout |
| `LOG_LEVEL` | `info` | Log level: debug, info, warn, error |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness; returns `{"status":"ok"}` |
| GET | `/process` | Process with default payload `"hello"` |
| POST | `/process` | Process body: `{"payload":"...", "echo": true\|false}` |

Response for `/process`: `{"result":"...", "processed_at_unix": <unix_ts>}`.

## Example requests

**Bash / curl:**

```bash
curl http://localhost:8080/health
curl http://localhost:8080/process
curl -X POST http://localhost:8080/process -H "Content-Type: application/json" -d '{"payload":"test", "echo": false}'
```

**PowerShell:**

```powershell
Invoke-WebRequest -Uri http://localhost:8080/health
Invoke-WebRequest -Uri http://localhost:8080/process
Invoke-WebRequest -Uri http://localhost:8080/process -Method POST -ContentType "application/json" -Body '{"payload":"test"}'
```

Stop the server with Ctrl+C or SIGTERM; it will shut down gracefully within `SHUTDOWN_TIMEOUT`.
