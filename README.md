# mini-api-golang

A lightweight RESTful API service built with Go, featuring user management and task management, backed by SQLite and secured with JWT authentication.

## Tech Stack

| Component      | Technology                          |
|---------------|-------------------------------------|
| Web Framework | [Gin](https://github.com/gin-gonic/gin) v1.9 |
| Database      | SQLite (via [GORM](https://gorm.io)) |
| Auth          | JWT (HS256, 72-hour expiry)         |
| Config        | [Viper](https://github.com/spf13/viper) (YAML) |
| Logging       | [Zap](https://github.com/uber-go/zap) |
| Go version    | 1.21+                               |

---

## Project Structure

```
mini-api-golang/
├── cmd/
│   └── main.go            # Application entry point
├── config/
│   ├── app.yaml           # Configuration file
│   └── config.go          # Config loader (Viper)
├── internal/
│   ├── dao/               # Data access layer (GORM models + queries)
│   │   ├── database.go    # SQLite initialization & auto-migration
│   │   ├── user_dao.go    # User CRUD operations
│   │   └── task_dao.go    # Task CRUD operations
│   ├── handler/           # HTTP handlers (Gin)
│   │   ├── user_handler.go
│   │   └── task_handler.go
│   ├── middleware/        # Gin middleware
│   │   ├── jwt.go         # JWT authentication middleware
│   │   └── logging.go     # Request logging middleware
│   ├── models/            # Domain models
│   │   └── user.go
│   ├── routes/
│   │   └── routes.go      # Route registration
│   ├── service/           # Business logic layer
│   └── utils/             # Shared utilities (JWT, response, password)
└── pkg/
    └── logger/            # Zap logger initializer
```

---

## Configuration

The application reads `config/app.yaml`. A sample configuration:

```yaml
server:
  port: 8080          # Port the HTTP server listens on

database:
  path: ./mini-api.db # Path to the SQLite database file

jwt:
  secret: change-me-in-production  # HMAC secret for JWT signing

log:
  level: info         # Log level: debug | info | warn | error
```

> **Security note:** Always replace `jwt.secret` with a strong, random value in production. Never commit secrets to version control.

When `log.level` is set to `debug`, the Gin engine runs in debug mode (verbose routing output). All other values switch it to release mode.

---

## Deployment

### Prerequisites

- Go **1.21** or later — [Download](https://go.dev/dl/)
- Git

### 1. Local Deployment (Development)

```bash
# 1. Clone the repository
git clone https://github.com/rockyfang2024/mini-api-golang.git
cd mini-api-golang

# 2. Download dependencies
go mod tidy

# 3. (Optional) Edit configuration
#    Open config/app.yaml and adjust port, database path, and jwt.secret

# 4. Build and run from the project root
go run ./cmd/main.go
```

The server will start and you will see log output similar to:

```
{"level":"info","msg":"configuration loaded","port":8080}
{"level":"info","msg":"database initialized","path":"./mini-api.db"}
{"level":"info","msg":"starting server","addr":":8080"}
```

The SQLite database file (`mini-api.db`) is created automatically on first run via GORM auto-migration.

### 2. Build a Binary

```bash
# Build for the current platform
go build -o mini-api ./cmd/main.go

# Run the binary
./mini-api
```

To cross-compile for Linux on macOS/Windows:

```bash
GOOS=linux GOARCH=amd64 go build -o mini-api-linux ./cmd/main.go
```

### 3. Docker Deployment

Create a `Dockerfile` in the project root:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mini-api ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/mini-api .
COPY config/app.yaml config/app.yaml
EXPOSE 8080
CMD ["./mini-api"]
```

Build and run:

```bash
docker build -t mini-api-golang .
docker run -d -p 8080:8080 --name mini-api mini-api-golang
```

To persist the SQLite database across container restarts, mount a volume:

```bash
docker run -d -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e MINI_API_DATABASE_PATH=/app/data/mini-api.db \
  --name mini-api mini-api-golang
```

### Graceful Shutdown

The server handles `SIGINT` and `SIGTERM` signals, waiting up to 5 seconds for in-flight requests to complete before exiting. Press `Ctrl+C` or send `kill <PID>` to trigger a clean shutdown.

---

## API Reference

Base URL: `http://localhost:8080`

### Response Format

All endpoints return a uniform JSON envelope:

```json
{
  "success": true,
  "message": "human-readable message",
  "data": { }
}
```

On error, `success` is `false` and `data` is omitted:

```json
{
  "success": false,
  "message": "error description"
}
```

---

### Health Check

#### `GET /health`

Returns the service status. No authentication required.

**Response `200 OK`:**
```json
{"status": "ok"}
```

---

### Authentication

#### `POST /register` — Register a new user

**Request body:**
```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "secret123"
}
```

| Field      | Type   | Required | Constraints        |
|------------|--------|----------|--------------------|
| `username` | string | ✅       | Non-empty          |
| `email`    | string | ✅       | Valid email format |
| `password` | string | ✅       | Minimum 6 chars    |

**Response `201 Created`:**
```json
{
  "success": true,
  "message": "user registered",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "created_at": "2026-03-14T04:00:00Z",
    "updated_at": "2026-03-14T04:00:00Z"
  }
}
```

**Response `409 Conflict`** — username or email already exists.

---

#### `POST /login` — Log in and obtain a JWT token

**Request body:**
```json
{
  "username": "alice",
  "password": "secret123"
}
```

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

The token is a **HS256-signed JWT** valid for **72 hours**. Use it as a Bearer token in the `Authorization` header for all protected endpoints:

```
Authorization: Bearer <token>
```

**Response `401 Unauthorized`** — invalid credentials.

---

### Users  *(requires authentication)*

All `/users` endpoints require an `Authorization: Bearer <token>` header.

#### `GET /users/:id` — Get a user by ID

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/users/1
```

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "created_at": "2026-03-14T04:00:00Z",
    "updated_at": "2026-03-14T04:00:00Z"
  }
}
```

**Response `404 Not Found`** — no user with that ID.

---

#### `PUT /users/:id` — Update a user's email

**Request body:**
```json
{
  "email": "newalice@example.com"
}
```

| Field   | Type   | Required | Constraints        |
|---------|--------|----------|--------------------|
| `email` | string | ❌       | Valid email format |

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "user updated",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "newalice@example.com",
    "created_at": "2026-03-14T04:00:00Z",
    "updated_at": "2026-03-14T04:10:00Z"
  }
}
```

---

#### `DELETE /users/:id` — Delete a user

```bash
curl -X DELETE -H "Authorization: Bearer <token>" http://localhost:8080/users/1
```

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "user deleted"
}
```

---

### Tasks  *(requires authentication)*

All `/tasks` endpoints require an `Authorization: Bearer <token>` header.

#### `POST /tasks` — Create a task

**Request body:**
```json
{
  "title": "Buy groceries"
}
```

| Field   | Type   | Required | Constraints |
|---------|--------|----------|-------------|
| `title` | string | ✅       | Non-empty   |

**Response `201 Created`:**
```json
{
  "success": true,
  "message": "task created",
  "data": {
    "ID": 1,
    "Title": "Buy groceries",
    "Done": false
  }
}
```

---

#### `GET /tasks/:id` — Get a task by ID

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/tasks/1
```

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "ID": 1,
    "Title": "Buy groceries",
    "Done": false
  }
}
```

**Response `404 Not Found`** — no task with that ID.

---

#### `PUT /tasks/:id` — Update a task

**Request body** (all fields optional):
```json
{
  "title": "Buy groceries and cook dinner",
  "done": true
}
```

| Field   | Type    | Required | Description             |
|---------|---------|----------|-------------------------|
| `title` | string  | ❌       | New task title          |
| `done`  | boolean | ❌       | Mark task as done/undone|

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "task updated",
  "data": {
    "ID": 1,
    "Title": "Buy groceries and cook dinner",
    "Done": true
  }
}
```

---

#### `DELETE /tasks/:id` — Delete a task

```bash
curl -X DELETE -H "Authorization: Bearer <token>" http://localhost:8080/tasks/1
```

**Response `200 OK`:**
```json
{
  "success": true,
  "message": "task deleted"
}
```

---

## End-to-End Usage Example

```bash
BASE=http://localhost:8080

# 1. Register
curl -s -X POST $BASE/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"secret123"}'

# 2. Login and capture the token
TOKEN=$(curl -s -X POST $BASE/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 3. Create a task
curl -s -X POST $BASE/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries"}'

# 4. Mark the task as done (replace 1 with the actual task ID)
curl -s -X PUT $BASE/tasks/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"done":true}'

# 5. Fetch the task
curl -s -H "Authorization: Bearer $TOKEN" $BASE/tasks/1
```

---

## Running Tests

```bash
go test ./...
```

---

## License

This project is open source. See the repository for license details.