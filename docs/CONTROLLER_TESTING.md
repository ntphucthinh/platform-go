# Controller Testing Guidelines

This document outlines the guidelines for implementing controller tests under `src/internal/controller_test`. The objective is to verify controller behavior via HTTP, ensuring integration testing with database persistence and end-to-end I/O coordination across the application.

## Directory Structure

```text
src/internal/controller_test/
  shared/            // Shared test utilities
    request.go       // Test user & request generation helpers
    router.go        // Gin router and external service mock setup
  auth_controller/   // Tests for authentication controller
    login_test.go    // Naming format: controller method name + _test.go
    logout_test.go
  other_controller/  // Tests for other controllers
    health_check_test.go
```

- Code in the `shared` directory is shared across all controller tests. Controller-specific logic should be placed under the corresponding `*_controller` directory.
- Directory names match the target controller (e.g., `project_controller`).
- Test file names under `*_controller` follow the format `<method_name>_test.go`, mirroring the controller method being tested.

## Using Shared Utilities

### Router Initialization

```go
r := shared.NewSetupRouter(t)
```

- After running overall test initialization via `utils.MainSetup()`, this sets up the Gin engine, test database, and external service mocks (e.g., S3 / ChatWork).
- The test GORM DB and external mocks can be accessed via the `Pkg` field of the returned `shared.Router`.

### Creating Test Users

```go
user := shared.CreateLoginUser(r.Pkg)
skUser := shared.CreateLoginUserWithSK(r.Pkg)
```

- Factory helpers insert required master data and user records into the test database.
- Related entities (such as projects, shops, etc.) are generated simultaneously.
- If specific roles are needed, create additional helper functions (refer to `CreateLoginUserWithSK` as a reference).

### Creating Requests

| Helper | Purpose |
| --- | --- |
| `Router.NewUnAuthRequest` | For public endpoints without authentication. Accepts HTTP method, URL, and optional payload body. |
| `Router.NewRequest` | For authenticated endpoints. Automatically logs in with a test user, attaching the CSRF token and cookies. |
| `Router.Login` | Used when explicit control of login state/cookies is needed within a test. |

- The request body is accepted as `interface{}` and serialized to JSON internally. Extend helpers if binary data handling is required.
- The default `Content-Type` is `application/json`. Add headers explicitly within the test when deviations are needed.

## Writing Tests

1. **Setup**: Initialize the environment using `shared.NewSetupRouter` and seed required test data using Factories.
2. **Create Request**: Choose `NewRequest` or `NewUnAuthRequest` depending on whether authentication is required.
3. **Assert Response**:
   - HTTP status code (e.g., `assert.Equal(t, http.StatusOK, w.Code)`)
   - JSON parsing (`json.Unmarshal` into generated entity types)
   - Cookies / headers / response messages
4. **Cover Negative Paths**: Explicitly cover branch conditions handled by the controller, such as authentication errors or validation failures.

### Example: Successful Login

```go
loginReq := api.AuthLoginRequest{LoginID: user.LoginID, Password: user.Password}
w := r.NewUnAuthRequest(t, "POST", "/auth/login", loginReq)

var res api.AuthLoginResponse
require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
assert.Equal(t, http.StatusOK, w.Code)
assert.NotEmpty(t, res.Data.CsrfToken)
```

## Running Tests with Database Connection

In controller tests, migrations and seeds are executed against the dedicated test MySQL database, and assertions typically run within transactions. Adjust router creation timing based on test scenarios to balance execution speed with test isolation.

### Setup Patterns

#### Reusing a Single Setup

```go
func TestController_Method(t *testing.T) {
	r := shared.NewSetupRouter(t)
	user := shared.CreateLoginUser(r.Pkg)

	t.Run("Normal_Success", func(t *testing.T) {
		w := r.NewRequest(t, user, "GET", "/uri", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var res api.AuthMeResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, user.ID, res.Data.ID)
	})
}
```

- Shares the router and test DB across all subtests.
- Faster execution because migration and seeding run only once, but requires care to avoid data pollution.
- If cross-test side effects are a concern, reset the DB directly or use the per-test reset pattern below.

#### Resetting Per Test Case

```go
func TestController_Method(t *testing.T) {
	t.Run("Normal_Success", func(t *testing.T) {
		r := shared.NewSetupRouter(t)
		user := shared.CreateLoginUser(r.Pkg)

		w := r.NewRequest(t, user, "GET", "/uri", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var res api.AuthMeResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, user.ID, res.Data.ID)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		r := shared.NewSetupRouter(t)
		user := shared.CreateLoginUser(r.Pkg)

		w := r.NewRequest(t, user, "GET", "/uri", nil)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Contains(t, "Authentication required", res.Message)
	})
}
```

- High test isolation because migrations and seeds run for each subtest.
- Execution takes longer, but ideal when error test cases deliberately mutate or pollute database state.

### Test Data Generation

- Factory functions generate unique identifiers using `time.Now().UnixNano()` to prevent unique constraint collisions.
- When an entity lacks a factory, add it under `tests/factory` and expose it through the `shared` package if appropriate.
- Key helpers:
  - `shared.CreateLoginUser` creates a standard user
  - `shared.CreateLoginUserWithSK` creates an SK user
  - `factory.NewXXX` functions provide `With` methods to override attributes

### Request Dispatch Helpers

| Helper | Role | Details |
| --- | --- | --- |
| `Router.NewUnAuthRequest` | Requests to public endpoints | Automatically sets `Content-Type: application/json` and serializes the body to JSON |
| `Router.NewRequest` | Requests to authenticated endpoints as a logged-in user | Calls `Router.Login` internally to attach CSRF headers and cookies |
| `Router.Login` | Acquires login session and cookies on demand | Useful for negative scenarios like tampered refresh tokens |

If custom headers are required, configure them using `req.Header.Set` after creating the `http.Request`.

### Response Assertion Patterns

- **HTTP Status**: Verify expected status codes like `http.StatusOK` or `http.StatusUnauthorized`.
- **JSON Parsing**: Unmarshal into `generated/api` models and assert with `require.NoError`.
- **Response Content**: Verify critical fields such as IDs or messages using `assert.Equal` / `assert.NotEmpty`.
- **Cookies & Headers**: Inspect `w.Result().Cookies()` to verify presence, attributes, or clearing of JWT, CSRF, and refresh tokens.
- **Database State**: For mutating requests, verify record existence and column values directly using `r.Pkg.TestDB`.

### Scenario Design Tips

- Structure tests clearly as `t.Run("Normal_...", func(t *testing.T) { ... })` and `t.Run("Error_...", func(t *testing.T) { ... })`.
- In negative paths, create individual cases for each controller branch (e.g., deleted/tampered cookies, non-existent login IDs).
- Use descriptive test names and assertion messages aligned with specifications so failures can be quickly diagnosed.

### Best Practices

- Ensure test reproducibility by generating unique data via factories.
- Use `require` for preconditions so tests fail fast before entering assertions.
- Keep response models and error messages synchronized with OpenAPI definitions; run `make bundle && make codegen` after spec changes.
- For tests with side effects, inspect the DB post-test and ensure rollback semantics are maintained.
- Validate calls to external services (S3, ChatWork, etc.) using existing mocks like `r.Pkg.TestS3Client`.

### Implementation Checklist

1. Reuse helpers under `shared` to eliminate duplicate setup code.
2. Verify both HTTP status codes and response bodies across happy and error paths.
3. For authenticated endpoints, always use `Router.NewRequest` to ensure CSRF headers and cookies are attached.
4. Verify database state transitions using `r.Pkg.TestDB` to ensure side effects are contained.
5. Use descriptive test names and failure messages tied to business specifications.

## Extension and Maintenance

- Place new controller tests following `<controller>_controller/<endpoint>_test.go`.
- Extract common logic into `shared` as patterns emerge.
- When replacing or configuring external service mocks, utilize `shared.Router.Pkg` mocks and configure mock expectations within the test.
- Maintain any newly created factories or mocks in accordance with existing patterns under `tests/factory` or `mocks`.

## Expected Test Scope

- Authentication (`auth_controller`): Login, logout, token refresh, current user info
- Health checks & system status (`other_controller`)
- Domain & feature controllers (`project_controller`, `user_controller`, etc.)

Base all test assertions on the OpenAPI contract returned by controllers: JSON schema, status codes, and observable side effects (cookies, DB mutations).
