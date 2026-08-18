# Todo REST API in Go

A REST API for authenticated todo management, built with Go, [Gin](https://gin-gonic.com/), [GORM](https://gorm.io/), and PostgreSQL. Passwords are hashed with bcrypt and login returns a JWT valid for 24 hours.

## Prerequisites

- Go 1.26.5 or later
- PostgreSQL running locally or remotely

## Configuration

Create a `.env` file in the project root, or provide these values through the system environment:

```dotenv
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=todo_app
JWT_SECRET=replace_with_a_long_random_secret
```

The application connects to PostgreSQL with `sslmode=disable`. On startup, GORM automatically creates or updates the `users` and `todos` tables.

## Running the API

```bash
go mod tidy
go run .
```

The server listens on `http://localhost:8080`.

## Authentication

Signup and login are public. Every `/todos` endpoint requires the token returned by login:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"yasmin","password":"test123"}'

curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"yasmin","password":"test123"}'
```

Set the returned token in the `Authorization` header:

```bash
TOKEN="paste_login_token_here"
curl http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN"
```

New accounts have the `user` role. The `DELETE /todos` endpoint is restricted to users whose role is `admin`.

## Endpoints

### Public

| Method | Path | Description |
|---|---|---|
| POST | `/signup` | Create a user account |
| POST | `/login` | Authenticate and return a JWT |

### Authenticated todo routes

All routes in this table require `Authorization: Bearer <token>`.

| Method | Path | Description |
|---|---|---|
| GET | `/todos` | Return all todos |
| GET | `/todos/:id` | Return one todo by ID |
| GET | `/todos/category/:category` | Filter todos by category |
| GET | `/todos/status/:status` | Filter by completion status (`true` or `false`) |
| GET | `/todos/search?q=term` | Case-insensitive title search |
| POST | `/todos` | Create a todo for the authenticated user |
| PUT | `/todos/:id` | Update a todo |
| PUT | `/todos/category/:category` | Update completion for every todo in a category |
| DELETE | `/todos/:id` | Delete a todo; regular users may delete only their own |
| DELETE | `/todos` | Delete all todos; admin only |

## Todo fields

Todo request and response JSON uses these fields:

| Field | Type | Notes |
|---|---|---|
| `id` | integer | Assigned by the database |
| `title` | string | Required |
| `completed` | boolean | Defaults to `false` |
| `category` | string | Defaults to `General` |
| `priority` | string | Must be `Low`, `Medium`, or `High` (case-insensitive input) |
| `completedAt` | timestamp or null | Set automatically when completed |
| `dueDate` | timestamp or null | Must be in the future when provided |
| `userID` | integer | Set from the authenticated user on creation |

Example create request:

```bash
curl -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Finish Go project","category":"Work","priority":"High","dueDate":"2030-01-01T12:00:00Z"}'
```

Example update request:

```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated title","completed":true,"category":"Work","priority":"High"}'
```

## Error responses

| Status | Cause |
|---|---|
| `400` | Malformed JSON, invalid ID/status, missing title, invalid priority, or past due date |
| `401` | Missing, invalid, expired, or malformed JWT |
| `403` | Insufficient permissions or deleting another user's todo |
| `404` | No todo found with the given ID or category |
| `409` | Username already exists |
| `500` | Database or server error |

## Testing

```bash
go test ./...
```

## Project structure

| File | Description |
|---|---|
| `main.go` | Application entry point and route registration |
| `connect.go` | PostgreSQL connection and GORM migrations |
| `auth.go` | Signup, login, password hashing, and JWT creation |
| `middleware.go` | JWT authentication and admin authorization |
| `todo.go` | Todo model |
| `user.go` | User model |
| `todo_handlers.go` | Todo request validation and HTTP handlers |
| `repository.go` | GORM todo persistence operations |
| `user_repository.go` | GORM user persistence operations |
| `repository_interface.go` | Repository contracts used by handlers and tests |
| `*_test.go` | Handler and authentication tests |