# Composer API

A set of APIs to allow users to manage their available features on their workspaces.

## Prerequisites

- Go 1.22+
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (for running migrations locally)
- [golangci-lint](https://golangci-lint.run/) (for linting)
- Podman or Docker (for container builds)

## Getting Started

```bash
cp configs/config.sample.yaml configs/config.yaml  # adjust values as needed
make run                                            # tidy, build, and start the server
```

The server listens on port `8080` by default.

## Build & Run

```bash
make          # build binary to bin/composer-api
make run      # tidy + build + run
make test     # run tests with -race and coverage
make lint     # golangci-lint
make image    # build container image (auto-detects podman/docker)
make container # run container on port 8080
make clean    # remove bin/ and coverage.out
```

## Database Migrations

Migrations live in `internal/database/migrations/` and are managed with [golang-migrate](https://github.com/golang-migrate/migrate).

```bash
# Create a new migration pair
make migrate-create name=<migration_name>

# Run all pending migrations (requires DATABASE_URL)
export DATABASE_URL="postgres://user:pass@localhost:5432/composer?sslmode=disable"
make migrate-up

# Rollback the last migration
make migrate-down
```

In Clowder deployments, migrations run automatically via an init container before the application starts. The init container reads database credentials from the Clowder config (`/cdapp/cdappconfig.json`) and executes all pending migrations.

## Configuration

Configuration is loaded via [Viper](https://github.com/spf13/viper) in the following order of precedence:

1. Clowder config (when deployed on Clowder, database settings are auto-populated)
2. `config.yaml` file (searched in `.`, `configs/`, `/etc/composer-api/`)
3. Built-in defaults

See `configs/config.sample.yaml` for all available options.

## Project Structure

```text
cmd/composer-api/        Entry point, graceful shutdown
configs/                 Sample and local configuration files
deploy/                  Clowder deployment (clowdapp.yaml, migrate.sh)
internal/
  config/                Viper-based config with Clowder integration
  database/
    migrations/          Sequentially numbered .up.sql / .down.sql files
  handler/
    probe/               Liveness/readiness endpoints (ping, health)
    v1/workspace/        Workspace API handlers
  logger/                Structured logging with request-scoped context
  middleware/            Request ID, logging, recovery, CORS, body size limit
  response/              JSON response helpers and typed HTTP errors
  server/                HTTP server and route registration
```

## Deployment

The application is deployed on OpenShift via [Clowder](https://github.com/RedHatInsights/clowder). The deployment manifest is at `deploy/clowdapp.yaml`.

The container image includes the `migrate` CLI. An init container runs `deploy/migrate.sh` to apply pending database migrations before the application starts.
