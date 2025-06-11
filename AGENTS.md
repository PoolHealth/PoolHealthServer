# Repository Guidelines

This repository contains the source for `PoolHealthServer`, a Go based server with a GraphQL API for pool management.  The project uses standard Go modules and follows common Go conventions.

## Layout
- `cmd/` contains program entry points like the main server.
- `internal/` holds application packages and services.
- `pkg/` exposes versioned API packages (GraphQL generated code lives in `pkg/api/v1`).
- `common/` defines shared types.
- `deployment/`, `Dockerfile`, and `docker-compose.yml` provide deployment resources.

## Development
1. **Formatting**: run `go fmt ./...` before committing any Go code.
2. **Dependencies**: if modules change, run `go mod tidy` and commit the updated `go.mod` and `go.sum` files.
3. **Testing**: run `go test ./...` to execute all tests.  Some tests rely on external credentials and may fail when these secrets are missing.
4. **Linting**: the project uses `golangci-lint` with `.golangci.yml`.  Run `golangci-lint run ./...` to check code style.
5. **GraphQL**: if you modify schema files under `pkg/api/v1/graphql`, regenerate code using `go generate ./...`.

## Pull Requests
- Write clear commit messages describing what changed and why.
- Ensure programmatic checks (tests and lint) pass when possible.  If they fail due to missing dependencies, note this in your PR description.
- Keep changes scoped and focused on a single purpose.
