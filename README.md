# Go Service Template

REST API service template with Chi framework. It provides configuration, PostgreSQL, a sample CRUD resource, validation, JSON errors, logging, graceful shutdown, and tests without introducing unnecessary framework complexity.

The sample domain is `Task`. Rename or replace `internal/task` with your resource (for example `users`, `orders`, or `books`) and keep the request flow intact.

## Quick start

1. Create your local environment file: `cp .env.example .env`.
2. Start PostgreSQL: `docker compose up -d postgres`.
3. Export the configuration: `set -a; source .env; set +a`.
4. Download modules: `go mod tidy`.
5. Start the API: `go run ./cmd/api`.

`GET /health` responds with `{"status":"running"}` after confirming that PostgreSQL is reachable.

## API

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `GET` | `/health` | Liveness plus database connectivity |
| `POST` | `/api/v1/tasks` | Create a task |
| `GET` | `/api/v1/tasks` | List tasks |
| `GET` | `/api/v1/tasks/{id}` | Get one task |
| `PUT` | `/api/v1/tasks/{id}` | Replace a task |
| `DELETE` | `/api/v1/tasks/{id}` | Delete a task |

Example:

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Design the schema","description":"Start with entities and relationships","status":"todo"}'
```

## How the code is organized

```text
cmd/api                 application wiring and graceful shutdown
internal/config         environment variables only
internal/database       PostgreSQL connection and pool settings
internal/httpapi        routes, middleware, common JSON responses
internal/task           model, repository, service, HTTP handler, tests
```

For a request such as `POST /api/v1/tasks`, the handler parses HTTP and JSON, the service validates business input, and the repository persists it. That separation gives you a clean place to add rules and makes service/handler tests possible without a real database.

## Architecture notes

- Define the entities, relationships, constraints, and API contract before implementing the smallest end-to-end vertical slice.
- Keep handlers thin. Put validation and business rules in the service, and SQL/ORM calls in the repository.
- Use `400` for malformed or invalid input, `404` for missing resources, `201` after creation, and `204` after deletion. Never return raw database errors to clients.
- `AutoMigrate` is convenient during early development. For a production service, replace it with versioned SQL migrations after the schema stabilizes.
- The included test uses an in-memory repository, so it runs quickly and demonstrates how to add more tests without PostgreSQL.

## Verification

Run `go test ./...`, start the database, and make one successful `curl` request.
