# AGENTS.md

## Languages & Frameworks

- Use **Go 1.24** (as specified in `src/go.mod`) as the primary language.
- Manage API specifications with **OpenAPI** (`src/openapi/openapi.yml`, split under `src/openapi/components/` and `src/openapi/paths/`) and generate code with **oapi-codegen**.
- Use **MySQL 8.0.41** as the primary database (containers `db` / `db-test`).
- Use **Docker Compose** (`docker-compose.yml`) to orchestrate local application, database, swagger, and tooling services.

## Architecture & Design Principles

### Clean Architecture & Domain-Driven Design (DDD)

This project follows **Clean Architecture** with **DDD** to keep business rules isolated from frameworks and infrastructure.

#### Layer Structure

- **Domain Layer** (`src/internal/domain/`):
  - Core business entities, value objects, and domain services.
  - Must remain free from transport, database, or framework dependencies.

- **UseCase Layer** (`src/internal/usecases/`):
  - Application-specific workflows orchestrating domain objects.
  - Keep modules focused (e.g., project, selection, information).

- **Interface Layer** (`src/internal/interfaces/controller/`):
  - HTTP controllers and adapters mapping requests/responses to use cases.
  - Organised by controller type (admin, dealer, project, etc.).

- **Infrastructure Layer** (`src/internal/infrastructure/`):
  - External integrations: repositories, query services, AWS, ChatWork, websocket, DB utilities.
  - Implement interfaces defined in the domain/usecase layers.

- **Supporting Packages**:
  - Shared cross-cutting code lives under `src/internal/common/`, `src/internal/middleware/`, and `src/internal/pkg/`.
  - Controller tests reside in `src/internal/controller_test/` and must stay aligned with interface-layer behaviour.

#### Dependency Injection Pattern

- Use **Google Wire** for compile-time dependency injection.
- Main graph: `src/internal/di/wire.go`; generated output: `src/internal/di/wire_gen.go`.
- Test graph: `src/internal/di/wire-test.go`; generated output: `src/internal/di/test_wire_gen.go`.
- Keep constructors deterministic and side-effect free to simplify wiring.

### Code Quality Principles

- **Readability first**: favour expressive names, cohesive functions, and limited side effects.
- **Performance aware**: optimise database access (batching, avoiding N+1), leverage query services for heavy reads, and consider caching when appropriate.
- **Maintainability**: respect SOLID principles, keep modules small, and document invariants when code is non-obvious.
- **Testability**: design units with explicit dependencies; prefer interfaces for infrastructure to enable mocking.

## Dependency & Tooling Conventions

- Manage Go dependencies via **Go modules** (`src/go.mod`, `src/go.sum`); run `go mod tidy` after adding or upgrading packages.
- Execute build/test/utility commands through the root `makefile`; additional helpers live in `scripts/makefile`.
- Maintain OpenAPI definitions in `src/openapi/`; consolidate specs with `make bundle`, then generate Go/TypeScript code with `make codegen` (uses `src/server.yaml`, `src/models.yaml`).
- Generated artifacts output to `src/generated/`; do not hand-edit anything under this directory.
- Store database migrations in `src/db/migrations/` (paired up/down SQL files).
- Maintain seed data in `src/db/seed/seed.go` and supporting packages under `src/db/seed/seeder/`.
- Docker images are defined beneath `dockerfiles/`; Go service configs in `dockerfiles/go/`, MySQL configs in `dockerfiles/mysql/`, swagger/node helpers in `dockerfiles/node/`.
- AWS CLI dev credentials/scripts live at `dockerfiles/go/conf/aws-cli-dev/`.
- Use the bundled **phpMyAdmin** service for local DB inspection (`http://localhost:8080`).
- Follow documentation references under `docs/`, especially `docs/GO_CODING_GUIDELINES.md`, `docs/OPENAPI_GUIDELINES.md`, and `docs/TESTING_SCOPE.md`.

## Coding Style Rules

- Comply with the **Go Coding Guidelines** (`docs/GO_CODING_GUIDELINES.md`).
- Indent Go code with **tabs**; keep files/directories in **snake_case**.
- Avoid manual edits to generated files within `src/generated/` (including OpenAPI, Wire, and mock outputs).
- Reuse shared error handling/logging utilities from `src/internal/common/`.
- Provide tests according to `docs/TESTING_SCOPE.md`; target meaningful coverage for controllers, use cases, and domain logic.

## Controller Test Implementation

- Refer to `docs/CONTROLLER_TESTING.md` for detailed guidelines and examples.
- Place controller tests under `src/internal/controller_test/<controller>_controller` and name files `<endpoint>_test.go` to mirror controller methods.
- Bootstrap each test with `shared.NewSetupRouter(t)` to initialise Gin, database fixtures, and external service mocks; reuse `shared.Router` helpers for login and requests.
- Create test users and related entities via the shared factory helpers (`shared.CreateLoginUser`, `shared.CreateLoginUserWithSK`, etc.) before issuing requests.
- Use `Router.NewRequest` for authenticated endpoints and `Router.NewUnAuthRequest` for public ones; keep payloads as structs so helpers handle JSON serialisation.
- Validate HTTP status codes, response bodies, cookies/headers, and observable side effects (DB changes, emitted tokens) against the OpenAPI contract; cover both success and error paths.

## Agent-specific Guidelines

### Code Generation & Build

- After creating or updating Go files, run `make lint test` before committing (add `make import` if goimports cleanup is required).
- Use `make bundle` followed by `make codegen` to refresh OpenAPI-derived code; do **not** modify oapi-codegen outputs manually.
- Use `make wire` to regenerate dependency injection code (keeps `wire_gen.go` / `test_wire_gen.go` in sync).
- Use `make mock` (mockery) to generate interface mocks under `src/mocks/`.
- Always create both up and down SQL files when adding migrations (`make migration name=<snake_case_name>`).
- Prefer `sqlmock` for database unit tests and place helpers under `src/mocks/` or `src/testonly/` as established.
- Leave a clear follow-up note when assumptions are uncertain so reviewers can revisit them.
- Minimise manual edits to generated outputs, docker automation, and build scripts.

### Architecture Implementation

- Enforce dependency direction: `Domain ← UseCases ← Interfaces ← Infrastructure`.
- Keep domain entities/policies pure; infrastructure must satisfy interfaces defined closer to the core.
- Delegate business rules to domain/services; keep use cases focused on orchestration and transaction boundaries.
- Apply repository/query-service patterns consistently for data access and reporting use cases.

### Performance & Optimisation

- Vet database queries for N+1 patterns; batch or eager-load where feasible.
- Use query services (`src/internal/infrastructure/queryservice/`) for complex read models.
- Limit `SELECT *`; fetch required fields and ensure appropriate indexes.
- Introduce caching layers cautiously, guarding consistency requirements.
- Implement pagination or streaming for endpoints returning large datasets.

### Code Quality & Readability

- Use `github.com/cockroachdb/errors` (already depended) for contextual error wrapping.
- Maintain ubiquitous language from domain experts across code, tests, and docs.
- Adhere to existing naming and package structures (
  e.g., keep controller-specific logic in matching subdirectories).
- Document tricky code paths with precise comments and link to specs or tickets when helpful.

## Test Commands Reference

**Execute all `make` commands from the repository root (`platform-go/`).**

**Core workflow:**

- `make setup` – configure git hooks, build containers, and start the stack.
- `make build` / `make build-nc` – build Docker images (optionally without cache).
- `make up` / `make up-log` / `make down` / `make restart` / `make ps` – control Docker Compose services.
- `make go` / `make go-log` / `make air` – interact with the Go application container (shell, logs, hot reload via Air).
- `make node` – enter the Node toolchain container.
- `make re-swagger` – restart swagger-ui to reload bundled specs.
- `make import` – run goimports across the workspace.

**Quality & testing:**

- `make lint` – run golangci-lint inside the Go container.
- `make test` – reset the test database, seed fixtures, and run `go test -cover ./...` once per package.
- `make test-controller controller=<name>` – run controller-focused tests (see `src/internal/controller_test/`).
- `make test-controller-func controller=<name> func=<TestName> [l=<count>]` – target a specific test function (optional repeat loop via `l`).
- `make test-cf c=<name> f=<TestName> [l=<count>]` – shorthand wrapper for the previous command.
- `make test-d` – run verbose tests across all packages.
- `make coverage` – collect coverage profile and summary.
- `make coverage-d grep=<pattern>` – collect coverage then filter results with `grep`.
- `make mock` – regenerate mocks with mockery.

**Code generation:**

- `make bundle` – aggregate modular OpenAPI files into `src/generated/openapi/openapi/openapi.yaml` using the Node toolchain.
- `make codegen` – run oapi-codegen for server/models and refresh TypeScript definitions (triggers `make bundle` automatically).
- `make codegen-ts` – copy generated TypeScript definitions to the sibling frontend (invoked by `make codegen`).
- `make wire` – regenerate Wire dependency injection code.
- `make go-b` / `make go-s` – build or run the Go binary (local debugging).

**Database & migrations:**

- `make migration name=<snake_case>` – scaffold paired SQL migration files under `src/db/migrations/`.
- `make migrate-up` / `make migrate-down` / `make migrate-down-one` – apply or roll back migrations against the local MySQL instance.
- `make migrate-refresh` / `make migrate-fresh` – rebuild the application database (drop + recreate before applying migrations).
- `make seeder` – run `src/db/seed/seed.go` against the local DB.
- `make test-migrate-up` / `make test-migrate-down` / `make test-migrate-refresh` / `make test-migrate-fresh` – manage the dedicated test database (`db-test` container).
- `make test-seeder` – seed the test database with fixtures.
- `make db` / `make db-test` / `make test-reset-db` – open MySQL shells or reset the test schema directly.

**Important notes:**

- Run tests and migrations via `make` from the host; avoid executing `go test` or `migrate` manually inside containers unless debugging.
- Swagger UI serves the bundled spec at `http://localhost:8081`; HTTPS proxy exposed on `https://localhost:8082`.
- Non-`make` Go commands should be run inside the Go container (`make go`) to ensure toolchain parity.

**Testing best practices:**

- ✅ Prefer `make test`, `make test-controller`, or `make test-controller-func` from the project root.
- ✅ Use `make coverage` / `make coverage-d grep=<pattern>` to review coverage and spot gaps.
- ❌ Avoid running `go test` directly inside containers via `make go`; rely on the scripted targets to ensure migrations and seeds run correctly.
- ✅ Regenerate code (`make codegen`) and re-run targeted tests when OpenAPI contracts change.
- ✅ Regenerate DI (`make wire`) before testing when constructors change.
