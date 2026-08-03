# CLAUDE.md

## Build & Run

```bash
make          # build binary to bin/composer-api
make test     # run tests with -race and coverage
make lint     # golangci-lint
make run      # tidy + build + run
make image    # build container image (auto-detects podman/docker)
make container # run container on port 8080
make clean    # remove bin/ and coverage.out
```

## Project Structure

```text
cmd/composer-api/        Entry point, graceful shutdown
internal/
  config/                Viper-based config with validation (server.port, log.level, log.pretty)
  logger/                slog logger with ContextHandler; JSON (production) or text (pretty mode)
  middleware/            One file per middleware (requestid, logging, recover, chain, cors, bodysize)
  handler/
    probe/               Liveness/readiness endpoints (ping, health)
    v1/workspace/        v1 workspace API handlers
  response/              JSON response helpers, typed HTTP errors
  server/                HTTP server, route registration
```

## Conventions

### Handler pattern
Handlers are methods on a per-resource `Handler` struct with a `New()` constructor. Dependencies (DB, cache) are fields on the struct, injected via `New()`.

```go
// internal/handler/v1/workspace/workspace.go
type Handler struct { DB *sql.DB }
func New(db *sql.DB) *Handler { return &Handler{DB: db} }

// internal/handler/v1/workspace/list.go
func (h *Handler) List(w http.ResponseWriter, r *http.Request) { ... }
```

### Naming
- The file defining a resource's `Handler` struct is named after the package: `workspace.go` in package `workspace`, `probe.go` in package `probe`.
- Import aliases are not used unless two packages collide in the same file.

### Routing
- Go 1.22+ `net/http.ServeMux` with method-pattern routing (`"GET /ping"`).
- Versioned routes use sub-mux + `http.StripPrefix` (see `internal/server/routes.go`).
- Per-resource middleware stacks via `middleware.Chain`.

### Response shapes
- Success: `{"body": {...}}`
- Error: `{"status": "Not Found", "reason": "workspace not found"}`
- Use typed error factories: `response.NotFound().WithReasonStr("...")`, not raw status codes.
- `response.WriteError(w, err)` writes the error with the correct HTTP status.

### Logging
- Use `slog.InfoContext`/`slog.ErrorContext` (not `slog.Info`/`slog.Error`) so request_id propagates automatically via ContextHandler.

## Git Workflow

- Default branch: `main`
- Branch naming: use JIRA ticket ID as branch name (e.g., `KATANA-395`)
- Commit messages: short, imperative, focused on "why"
- PRs: title format `<JIRA-ID>: <description>`, include JIRA link in body
- Rebase on `main` before merging
