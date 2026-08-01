# Todo REST API in Go

A simple REST API for managing todo items, built with Go, [Gin](https://gin-gonic.com/), and [GORM](https://gorm.io/). Todos are persisted in a PostgreSQL database — data survives server restarts.

---

## Project structure

| File | Description |
|---|---|
| `main.go` | Entry point. Initializes the DB connection, repository, handlers, and registers all 5 routes. |
| `connect.go` | Database connection using GORM and environment variables.
| `models.go` | The `Todo` struct with GORM and JSON tags. |
| `repository.go` | All database logic — GORM-based implementation of Create, GetAll, GetByID, Update, Delete. |
| `handlers.go` | All route handler functions. Validates requests and delegates to the repository. |

---

## Prerequisites

- PostgreSQL running locally
- A `.env` file in the project root (see below)

---

## Environment variables

Create a `.env` file in the project root:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=todo_app
```

---

## Running the server

```bash
go mod tidy
go run .
```

The server starts at `http://localhost:8080`. The `todos` table is created automatically via GORM's `AutoMigrate` on startup.

---

## Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/todos` | Return all todos |
| GET | `/todos/:id` | Return a single todo by ID |
| POST | `/todos` | Create a new todo |
| PUT | `/todos/:id` | Update a todo's title and completion status |
| DELETE | `/todos/:id` | Delete a todo |

---

## Testing with cURL

### Create a todo
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Finish Go project"}'
```

### Get all todos
```bash
curl http://localhost:8080/todos
```

### Get a todo by ID
```bash
curl http://localhost:8080/todos/1
```

### Update a todo
```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated title", "completed": true}'
```

### Delete a todo
```bash
curl -X DELETE http://localhost:8080/todos/1
```

---

## Error responses

| Status | Cause |
|---|---|
| `400` | Malformed JSON, invalid ID, or empty title |
| `404` | No todo found with the given ID |
| `500` | Database error |

---