# Book API (Go)

REST API with user registration/login (JWT) and CRUD for books.
Go 1.22 standard-library router, PostgreSQL via pgx, bcrypt passwords.

## Quick start

```bash
cp .env.example .env          # then set JWT_SECRET (openssl rand -hex 32)
go mod tidy                   # downloads deps and creates go.sum (run once)
make db                       # start Postgres in Docker
make migrate-up               # create tables
make run                      # API on http://localhost:18743, docs on /docs
```

Or run everything in Docker: `docker compose up --build`

## Endpoints

| Method | Path                  | Auth          | Description                         |
|--------|-----------------------|---------------|-------------------------------------|
| POST   | /api/auth/register    | –             | Create account, returns token       |
| POST   | /api/auth/login       | –             | Log in, returns token               |
| GET    | /api/auth/me          | Bearer token  | Current user                        |
| GET    | /api/books            | –             | List (`?search=&limit=&offset=`)    |
| GET    | /api/books/{id}       | –             | Get one book                        |
| POST   | /api/books            | Bearer token  | Create book                         |
| PUT    | /api/books/{id}       | Bearer token  | Update book (owner only)            |
| DELETE | /api/books/{id}       | Bearer token  | Delete book (owner only)            |
| GET    | /healthz, /readyz     | –             | Liveness / readiness probes         |
| GET    | /docs                 | –             | Swagger UI (spec at /openapi.yaml)  |

## Try it

```bash
# register
curl -X POST localhost:18743/api/auth/register \
  -d '{"name":"Asha","email":"asha@example.com","password":"supersecret"}'

# login
TOKEN=$(curl -s -X POST localhost:18743/api/auth/login \
  -d '{"email":"asha@example.com","password":"supersecret"}' | jq -r .token)

# create
curl -X POST localhost:18743/api/books -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"The Go Programming Language","author":"Donovan & Kernighan","isbn":"9780134190440","published_year":2015}'

# list / get / update / delete
curl "localhost:18743/api/books?search=go"
curl localhost:18743/api/books/1
curl -X PUT localhost:18743/api/books/1 -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"The Go Programming Language","author":"Donovan & Kernighan","published_year":2016}'
curl -X DELETE localhost:18743/api/books/1 -H "Authorization: Bearer $TOKEN"
```

Errors come back as `{"error": "..."}`; validation errors are 422 with a `fields` map.

## Adding a new feature (like a new Django app)

1. Create `internal/<name>/` with `model.go`, `repository.go`, `service.go`, `handler.go`
2. Add a migration pair in `migrations/` (`000003_<name>.up.sql` / `.down.sql`)
3. Register routes in `internal/server/router.go`
4. Describe the new endpoints in `internal/docs/openapi.yaml`
