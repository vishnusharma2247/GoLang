# HTTP Learning Project (Golang + PostgreSQL + Docker)

This file is your reusable project blueprint.

From now onward, after each implementation step, I will append/update this README automatically with:
- one new substep under the correct big step
- commands used
- files touched
- purpose of the step

## Current Project Structure

```text
http/
├─ cmd/
│  └─ api/
│     └─ main.go
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  ├─ domain/
│  │  └─ user.go
│  ├─ repository/
│  │  └─ postgres/
│  │     ├─ db.go
│  │     └─ user_repository.go
│  ├─ security/
│  │  └─ password.go
│  ├─ service/
│  │  └─ auth_service.go
│  └─ transport/
│     └─ http/
│        └─ router.go
├─ migrations/
│  ├─ 000001_create_users.up.sql
│  └─ 000001_create_users.down.sql
├─ .env
├─ .env.example
├─ go.mod
├─ go.sum
└─ README.md
```

## Implementation Journey (Big Steps -> Substeps)

<details>
<summary><strong>Big Step 1: Project Bootstrap</strong></summary>

### Substep 1.1 - Initialize Go module
- Command: `go mod init http-learning`
- Output: `go.mod` created
- Why: starts module/dependency management

### Substep 1.2 - Create entrypoint
- File: `cmd/api/main.go`
- Added minimal HTTP server startup
- Why: app needs a runnable entrypoint first

### Substep 1.3 - Add first endpoint
- File: `cmd/api/main.go`
- Added `GET /health`
- Response: `200` + `{"status":"ok"}`
- Why: first HTTP sanity endpoint

</details>

<details>
<summary><strong>Big Step 2: Configuration Layer</strong></summary>

### Substep 2.1 - Add env template files
- Files: `.env.example`, `.env`
- Keys added:
  - `APP_ENV=development`
  - `APP_ADDR=:8080`
- Why: move runtime values outside code

### Substep 2.2 - Create config package
- File: `internal/config/config.go`
- Added:
  - `Config` struct
  - `Load()` function
  - `getEnv()` helper
- Why: centralized config access

### Substep 2.3 - Auto-load `.env`
- File: `cmd/api/main.go`
- Added: `godotenv.Load()`
- Commands:
  - `go get github.com/joho/godotenv`
  - `go mod tidy`
- Why: read `.env` automatically at app start

### Substep 2.4 - Add DB URL to env + config
- Files: `.env`, `.env.example`, `internal/config/config.go`
- Added key:
  - `DATABASE_URL=postgres://postgres:postgres@localhost:5432/http_learning?sslmode=disable`
- Added `DatabaseURL` in `Config`
- Why: prepare DB integration

</details>

<details>
<summary><strong>Big Step 3: HTTP Transport Refactor</strong></summary>

### Substep 3.1 - Create router package
- File: `internal/transport/http/router.go`
- Added `NewMux()` and moved `/health` route there
- Why: keep routing out of `main.go`

### Substep 3.2 - Wire router into main
- File: `cmd/api/main.go`
- Replaced inline route creation with `httptransport.NewMux()`
- Why: separation of concerns (`main` startup vs route definitions)

</details>

<details>
<summary><strong>Big Step 4: Domain + Database Foundation</strong></summary>

### Substep 4.1 - Add domain model
- File: `internal/domain/user.go`
- Added `User` struct:
  - `ID`, `Email`, `PasswordHash`, `Role`, `CreatedAt`
- Why: shared business entity for auth/repo/service

### Substep 4.2 - Add migrations
- Files:
  - `migrations/000001_create_users.up.sql`
  - `migrations/000001_create_users.down.sql`
- Up migration:
  - enable `pgcrypto`
  - create `users` table
- Down migration:
  - drop `users` table
- Why: version-controlled schema changes

### Substep 4.3 - Create DB pool initializer
- File: `internal/repository/postgres/db.go`
- Added `NewPool(databaseURL)`:
  - create pool
  - ping with timeout
  - return pool or error
- Commands:
  - `go get github.com/jackc/pgx/v5/pgxpool`
  - `go mod tidy`
- Why: robust DB connection management

### Substep 4.4 - Wire DB pool in startup
- File: `cmd/api/main.go`
- Added:
  - `postgres.NewPool(cfg.DatabaseURL)`
  - fail-fast error handling
  - `defer dbPool.Close()`
- Why: app should run only if DB is reachable

</details>

<details>
<summary><strong>Big Step 5: Real PostgreSQL via Docker</strong></summary>

### Substep 5.1 - Run PostgreSQL container
- Container name: `http-learning-db`
- Command:
```bash
docker run -d \
  --name http-learning-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=http_learning \
  -p 5432:5432 \
  postgres:latest
```
- Why: use real database in local dev

### Substep 5.2 - Validate DB container and SQL access
- Commands:
  - `docker ps --filter name=http-learning-db`
  - `docker logs --tail 20 http-learning-db`
  - `docker exec http-learning-db psql -U postgres -d http_learning -c "SELECT now();"`
- Why: confirm DB is healthy and queryable

### Substep 5.3 - Apply migration to DB
- Command:
```bash
docker exec -i http-learning-db psql -U postgres -d http_learning < migrations/000001_create_users.up.sql
```
- Verification:
```bash
docker exec http-learning-db psql -U postgres -d http_learning -c "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_name='users';"
```
- Why: create actual schema before repository calls

</details>

<details>
<summary><strong>Big Step 6: Repository Layer</strong></summary>

### Substep 6.1 - Add repository skeleton
- File: `internal/repository/postgres/user_repository.go`
- Added:
  - `UserRepository` struct
  - `NewUserRepository(pool)`
- Why: container for user-related SQL logic

### Substep 6.2 - Add create user query
- File: `internal/repository/postgres/user_repository.go`
- Added method:
  - `CreateUser(ctx, user)`
- SQL: `INSERT ... RETURNING id, created_at`
- Why: persist newly registered users

### Substep 6.3 - Add get user by email query
- File: `internal/repository/postgres/user_repository.go`
- Added method:
  - `GetUserByEmail(ctx, email)`
- Handles:
  - found user
  - not found (`nil, nil`)
  - DB error
- Why: login and duplicate email checks

</details>

<details>
<summary><strong>Big Step 7: Security Helpers</strong></summary>

### Substep 7.1 - Add password hash/check functions
- File: `internal/security/password.go`
- Added:
  - `HashPassword(plainPassword)`
  - `CheckPassword(hashedPassword, plainPassword)`
- Uses `bcrypt`
- Command: `go mod tidy`
- Why: never store plaintext passwords

</details>

<details>
<summary><strong>Big Step 8: Auth Service Logic</strong></summary>

### Substep 8.1 - Add service skeleton
- File: `internal/service/auth_service.go`
- Added:
  - `AuthService` struct
  - `NewAuthService(userRepo)`
- Why: service layer for auth business rules

### Substep 8.2 - Implement register flow
- File: `internal/service/auth_service.go`
- Added:
  - `Register(ctx, email, plainPassword)`
- Flow:
  - normalize input
  - validate
  - check existing user
  - hash password
  - create user via repository
- Added errors:
  - `ErrInvalidInput`
  - `ErrEmailAlreadyExists`
- Why: complete registration business logic

### Substep 8.3 - Implement login flow
- File: `internal/service/auth_service.go`
- Added:
  - `Login(ctx, email, plainPassword)`
- Flow:
  - normalize input
  - find user by email
  - compare bcrypt password hash
- Added error:
  - `ErrInvalidCredentials`
- Why: complete login credential validation logic

</details>

## How App Connects to Docker DB

1. `.env` has `DATABASE_URL` pointing to `localhost:5432`
2. `config.Load()` reads it
3. `main.go` calls `postgres.NewPool(cfg.DatabaseURL)`
4. Docker port mapping forwards host `5432` to container `5432`
5. Postgres in container authenticates and pool ping succeeds

Mental model:
- app runs on host
- database runs in container
- communication via mapped port

## Stable Port Rule

- App port: `:8080`
- DB port: `:5432`

If app port `8080` is busy:
- kill existing process using `8080`
- restart app on `:8080` (no temporary app ports)

## What Is Pending (Next Big Steps)

- HTTP auth handlers (`/register`, `/login`)
- JWT token generation + auth middleware
- Role-based authorization routes
- HTTPS local TLS setup
- Session-based auth flow
- Postman collection and testing scripts
