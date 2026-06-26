# Go HTTP Server

A small REST API written in Go using the standard `net/http` package, SQLite for
persistence, YAML/environment based configuration, and graceful shutdown handling.

The current API manages users:

- `POST /api/users` creates a user.
- `GET /api/users` returns all users.
- `GET /api/users/{id}` returns a single user by ID.

## Requirements

- Go `1.26.2` or newer, matching `go.mod`
- CGO-enabled Go toolchain for `github.com/mattn/go-sqlite3`
- SQLite-compatible local filesystem path for the database file

## Project Layout

```text
.
├── cmd/go-http-server/main.go          # Application entrypoint
├── config/local.yaml                   # Local configuration
├── internal/config                     # Config loading
├── internal/http/handlers/user         # User HTTP handlers
├── internal/storage                    # Storage interface
├── internal/storage/sqlite             # SQLite implementation
├── internal/types                      # Shared request/response types
└── internal/utils/response             # JSON response helpers
```

## Configuration

The server requires a configuration file. The path can be provided with either:

```sh
CONFIG_PATH=config/local.yaml go run ./cmd/go-http-server
```

or:

```sh
go run ./cmd/go-http-server -config config/local.yaml
```

Example configuration:

```yaml
env: "dev"
storage_path: "storage/storage.db"
http_server:
  address: "localhost:8000"
```

Environment variables can override configured values:

| Variable | Description | Default |
| --- | --- | --- |
| `APP_ENV` | Application environment name | `dev` |
| `APP_STORAGE_PATH` | SQLite database path | Required |
| `APP_HTTP_ADDRESS` | HTTP listen address | `:8000` |

## Run Locally

Install dependencies:

```sh
go mod download
```

Start the server:

```sh
CONFIG_PATH=config/local.yaml go run ./cmd/go-http-server
```

The local config starts the API at:

```text
http://localhost:8000
```

The SQLite storage layer creates the `users` table automatically if it does not
already exist.

## API

### Create User

```sh
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Aman","email":"aman@example.com","age":25}'
```

Successful response:

```json
{
  "success": true,
  "message": "user created successfully",
  "data": {
    "lastId": 1
  }
}
```

Required request fields:

| Field | Type |
| --- | --- |
| `name` | string |
| `email` | string |
| `age` | number |

### Get All Users

```sh
curl http://localhost:8000/api/users
```

Successful response:

```json
[
  {
    "id": 1,
    "name": "Aman",
    "email": "aman@example.com",
    "age": 25
  }
]
```

### Get User By ID

```sh
curl http://localhost:8000/api/users/1
```

Successful response:

```json
{
  "success": true,
  "message": "user retrieved successfully",
  "data": {
    "id": 1,
    "name": "Aman",
    "email": "aman@example.com",
    "age": 25
  }
}
```

## Development

Run all package tests:

```sh
go test ./...
```

Build the server binary:

```sh
go build ./cmd/go-http-server
```

## Notes

- User emails are stored with a database-level `UNIQUE` constraint.
- The server listens for `SIGINT` and `SIGTERM` and allows up to 5 seconds for
  graceful shutdown.
- The response helper wraps most success and error responses as JSON objects with
  `success`, `message`, and optional `data` fields.
