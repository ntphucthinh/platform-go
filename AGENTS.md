# AGENTS.md

## Project Overview

This repository is a Dockerized Go backend. The Go module and all application source are under `src/` (`module platform-go`, Go 1.24). It exposes a Gin HTTP API, persists to MySQL through GORM, and uses cookie-based JWT authentication. The implemented product scope is currently users plus administrator login; do not assume the broader modules described by older documentation exist.

Run Go commands from `src/`. Root-level `Makefile` commands orchestrate Docker services and code generation.

## Architecture

The code uses a small layered/CQRS-style arrangement rather than a complete generic Clean Architecture implementation:

```text
HTTP request
  -> Gin global middleware
  -> generated OpenAPI route wrapper
  -> interfaces/controller
     -> queryservice + qsdto                 (simple reads)
     -> usecase -> domain repository interface -> repository -> dbmodel -> MySQL (writes/business flows)
```

Authentication is a separate implemented vertical flow:

```text
middleware/controller -> AuthUsecase -> domain user.UserRepository + security.JWTService
```

`cmd/server/main.go` manually constructs the auth dependencies and uses Wire-generated DI for the user controller. Keep this split unless the application wiring is deliberately changed.

## Directory Structure

- `src/cmd/server/`: executable entry point; loads local env, opens GORM/MySQL, registers Gin middleware/routes, and starts the server.
- `src/internal/interfaces/controller/`: Gin handlers implementing `generated/api.ServerInterface`; maps HTTP/generated models to usecase or query-service calls and HTTP responses.
- `src/internal/middleware/`: global JWT and CSRF middleware.
- `src/internal/usecases/`: concrete application services and business/auth flows.
- `src/internal/domain/<feature>/`: dependency-light entities and repository contracts. At present only `domain/user` exists.
- `src/internal/infrastructure/repository/`: GORM implementations of domain repository interfaces, primarily command/write access.
- `src/internal/infrastructure/queryservice/` and `qsdto/`: direct read queries and their result DTOs. The current user list query goes straight from controller to query service.
- `src/internal/infrastructure/dbmodel/`: GORM table models and conversion to domain values.
- `src/internal/infrastructure/security/`: JWT signing/parsing and CSRF-token generation.
- `src/internal/common/`: HTTP controller helpers; `ctxkeys/` owns shared Gin-context key strings.
- `src/internal/pkg/app/`: application-environment helper; `pkg/types/` contains the authenticated-principal type. Do not turn these into business-logic catch-alls.
- `src/internal/di/`: Google Wire input and generated injectors.
- `src/db/migrations/`, `src/db/seed/`: SQL schema migrations and the Go seed executable.
- `src/openapi/openapi/`: hand-maintained, split OpenAPI source.
- `src/generated/`: generated OpenAPI artifacts.

Some empty directories (`internal/constants`, `infrastructure/auth`, `infrastructure/dbutil`, `mocks`, `testonly`, `tests`, and `openapi/paths`) are present but have no implementation. Do not infer a convention or add a parallel feature structure from them.

## Application Flow

`main` loads `.env.local` only when `APP_ENV` is empty or `local`, resolves `DB_DSN` or individual `DB_*` values, opens GORM, creates JWT/auth services, then registers CORS, `JWTAuth`, `CSRFProtection`, `/health`, and the generated OpenAPI handlers.

The JWT middleware skips `/login`, `/refresh`, and `/health`; it stores `types.LoginAdministrator` in Gin context under `ctxkeys.LoginAdministratorKey` (`"login_administrator"`). Authenticated endpoints also require the `csrf` cookie to match `X-CSRF-Token`, except `OPTIONS`, `/login`, and `/health`. Reuse that key/type rather than adding raw context strings.

## Domain, Usecase, and Controller Layers

Domain structs use plain Go fields with no persistence tags. Put a new feature's entity and repository interface together in `internal/domain/<feature>/` when it needs command-side persistence. Repository constructors return the domain interface (for example, `NewMySQLUserRepository(...) user.UserRepository`).

Usecases are concrete `*<Feature>Usecase` structs constructed by `New<Feature>Usecase`; they receive domain interfaces and accept `context.Context`. Put business decisions and multi-step application flows here. There are no usecase interfaces in the current codebase.

Controllers receive concrete usecase/query-service pointers. They bind and make basic request checks, call a usecase for commands or `QueryService` directly for simple reads, map values to generated `api` models, and write Gin responses. Do not call repositories or GORM directly from controllers. Generated `api` pointer fields are populated using local variables, as in `UserController`.

## Infrastructure and Database

Use GORM with MySQL and propagate the request context through `db.WithContext(ctx)`. `dbmodel` structs carry GORM tags and convert to domain values through `ToDomain()`. Query DTOs use explicit `gorm:"column:..."` tags and are populated by `Table(...).Find(...)`.

The repository does not currently centralize transaction handling; no explicit transaction pattern exists in application code. For a new atomic multi-write flow, keep the boundary explicit and close to the repository/usecase that owns it, following the surrounding implementation.

Migrations use `golang-migrate`, named `<timestamp>_<description>.up.sql` and matching `.down.sql`. Supply both directions when adding a migration. The migration history was consolidated into a single `create_users_table` migration before staging was first deployed; from now on, do not rewrite applied migration history—add a new migration instead. The seed program is `db/seed/seed.go` and is run inside the API container.

## Dependency Injection

Wire input files `internal/di/wire.go` and `wire-test.go` have the `wireinject` build tag. `wire_gen.go` and `test_wire_gen.go` are generated and compiled normally. Register new constructor providers in the input file, then run `make wire`; do not edit generated injector files. `InitializeHandler(db)` currently assembles the user repository, user usecase, user query service, and user controller.

## OpenAPI & Generated Code

The OpenAPI source of truth is `src/openapi/openapi/openapi.yml`, which references `paths/*.yml` and `components/schemas/*.yml`. Schema property spelling is lower camel case (for example `createdAt`, `accountType`); generated Go models use exported fields such as `CreatedAt` and `AccountType`.

For an API contract change:

1. Update the source OpenAPI root/path/schema files under `src/openapi/openapi/`.
2. Run `make codegen` (or `make gen`) from the repository root. It bundles the specification to `src/generated/openapi/openapi/openapi.yaml`, then invokes `oapi-codegen` using `src/config/server.yaml` and `src/config/models.yaml`.
3. Implement the newly generated `api.ServerInterface` method(s), update handler composition/wiring as needed, and test the behavior.

`src/generated/api/server.gen.go`, `src/generated/api/models.gen.go`, and the bundled `src/generated/openapi/` contents are generated outputs. Do not edit them directly. `platform-openapitools.json` is the active bundling configuration used by the Makefile. `src/config/openapitools.json` describes a TypeScript Axios generator but is not invoked by current Make targets.

## Error Handling and Logging

Errors are returned normally from repository/usecase/security layers. Existing sentinel error: `usecases.ErrInvalidCredentials`; `AuthUsecase.CurrentUser` uses `errors.Is` for `gorm.ErrRecordNotFound`. The security package uses `errors.New`; the one contextual wrap is `fmt.Errorf("generate CSRF token: %w", err)`.

At HTTP boundaries, existing code writes Gin JSON responses directly or uses `common.BindJSON`, `IsErrorCustom`, and `IsErrorWithMessage`. The common helpers return `{error, details}` for bind/internal failures; auth handlers intentionally return terse authentication errors. Follow the local endpoint's established response behavior rather than introducing a new global error system. Process startup/seed logging uses the standard `log` package; there is no project logging package or structured-logging convention.

## Testing

The only current tests are standard-library unit tests in `internal/infrastructure/security/jwt_test.go`, named `Test<Subject><Behavior>` and using direct `testing.T` assertions. There are no controller, integration, factory, mock, table-driven, or test-database tests yet, and no mock generation configuration. Keep new tests alongside the package under test unless a real integration-test structure is introduced. Prefer the standard `testing` style already present; use mocks only when a new test needs an interface dependency.

## Coding Conventions

- Use lowercase package/directory names and snake_case Go filenames (`user_usecase.go`, `jwt_middleware.go`). Exported Go names are PascalCase; unexported names are camelCase.
- Constructors are `New<Type>` and return pointers for controllers/usecases/services; the repository constructor returns its domain interface.
- Interfaces currently live with the domain model and are named by capability/entity (`UserRepository`), not `I...`.
- Use `context.Context` as the first method parameter for database/application operations and pass `c.Request.Context()` from handlers.
- Domain values have no tags. GORM models use `gorm` tags; OpenAPI controls generated JSON tags. The local authenticated principal uses explicit JSON tags.
- Source has sparse comments, mostly explaining exported symbols or non-obvious security/flow choices. Comment only when the intent is not clear from code.

## Adding a New Feature

1. Define or revise the OpenAPI contract under `src/openapi/openapi/` when an HTTP contract changes, then regenerate it.
2. Add the domain entity and command-side repository contract under `internal/domain/<feature>/` when persistence is needed.
3. Add a migration first for schema changes, then a GORM `dbmodel` and repository implementation. Add a `qsdto` plus query service for simple read/list endpoints.
4. Add a concrete usecase for commands/business flow. Keep simple reads in a query service only when they are just data retrieval, as with users.
5. Add the controller method(s), map to generated models, compose it into `Handler` if required, and update Wire input providers. Regenerate Wire.
6. Add focused tests, then format and run relevant checks.

## Common Commands

Run these from the repository root unless noted:

- `make build` / `make build-nc` — build Docker services.
- `make up`, `make down`, `make restart`, `make ps` — manage Docker Compose services. `make up` also configures repository Git hooks.
- `make air` — run Air hot reload in the `platform-api` container.
- `make lint` — run `golangci-lint` in the API container (no committed `.golangci.yml` is present, so it uses the tool's available defaults).
- `make tools` — install codegen/dev tools in the current Go environment.
- `make codegen` or `make gen` — bundle and generate OpenAPI outputs.
- `make wire` — regenerate Wire injectors.
- `make migration name=<migration_name>` — create a migration in the API container.
- `make migrate-up`, `make migrate-down`, `make migrate-refresh`, `make migrate-fresh`, `make seeder` — operate migrations/database/seed data. `migrate-down` drops all applied migrations; `migrate-fresh` drops and recreates `platform_db`, so use them only deliberately.
- From `src/`: `go test ./...`, `go test -race ./...`, `go build ./cmd/server`, and `go fmt ./...` are applicable standard Go commands. No repository `make test` or dedicated formatter target exists.

## Environment Configuration

`src/.env.local` is ignored by `src/.gitignore` and is loaded only for local app runs. It contains application, API database, and MySQL-container settings. `APP_ENV`, `PORT`, `DB_DSN` or `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_PASS`/`DB_NAME`, `JWT_SECRET`, and `CORS_ALLOW_ORIGIN` are used by the running server. Local fallback values exist for database fields; `JWT_SECRET` is required and the server fails fast when it is missing or shorter than 32 characters.

`src/.sample.env` is tracked but currently empty. `src/.env.test` is present locally but ignored; no `.sample.env.test` exists. Never place real secrets in this document or commit ignored `.env*` files. The root `.gitignore` also ignores `node_modules/`, `dist/`, `.DS_Store`, and `tmp/`; `src/.gitignore` ignores `.env*` and `tmp/`.

## Do Not Modify Directly

- Generated files under `src/generated/api/`, `src/generated/openapi/`, and Wire outputs `src/internal/di/wire_gen.go` and `test_wire_gen.go`; change their sources and regenerate.
- Local environment files `src/.env.local` and `src/.env.test`; they are ignored and may contain local credentials.
- Runtime temporary/log directories `src/tmp/` and `src/logs/`.
- Applied migration files: add a new migration rather than changing historical schema steps.

## Important Rules

- Work from actual source under `src/`; do not introduce removed application variants or parallel product-directory structures.
- Preserve the command/read split: controller -> usecase -> domain repository -> infrastructure for commands; controller -> query service for simple reads.
- Keep controller business logic and direct database access out of controllers; keep HTTP/Gin dependencies out of domain code.
- Use existing constructors, context propagation, domain interfaces, GORM models, and generated API types before adding abstractions or dependencies.
- Update OpenAPI and regenerate when an endpoint or schema changes; never hand-edit generated code.
- Keep JWT/CSRF middleware behavior and the shared authenticated-context key consistent for protected endpoints.
- Add or update tests for behavior changes, run formatting and relevant tests, and avoid unrelated edits.

## Known Documentation Inconsistencies

Root `README.md` is untracked in the current worktree and describes an older, much larger structure (for example `cmd/main`, `internal/controller`, multiple domains, and `.golangci.yml`). The implementation instead uses `cmd/server`, `internal/interfaces/controller`, one `user` domain, and has no committed lint config. The documents under `docs/` also contain outdated architecture references and commands such as `make generate`/`make dev` that do not exist. This file follows the checked source and Makefile instead.
