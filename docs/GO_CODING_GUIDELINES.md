# Go Coding Guidelines

## Basic Rules

| Category | Rule | Example | Notes |
| --- | --- | --- | --- |
| Indentation | Use tabs | `func sample() {`<br>&nbsp;&nbsp;&nbsp;&nbsp;`statements`<br>`}` | Spaces are discouraged |
| Line length | No limit | - | Break lines at appropriate positions when long |
| File naming | snake_case | `user_handler.go`<br>`db_connection.go` | All lowercase |
| File organization | 1 file = 1 package | `package main`<br>`package handlers` | - |
| Code layout order | 1. Type definitions (struct, interface)<br>2. Methods of those types<br>3. Related helper functions | - | - |

## Naming Conventions

| Element | Rule | Good Example | Bad Example |
| --- | --- | --- | --- |
| Package name | Single lowercase word | `http`<br>`strings` | `stringUtil`<br>`string_util` |
| Variable name | camelCase | `userID`<br>`firstName` | `UserID`<br>`first_name` |
| Constant name (Exported) | PascalCase | `MaxRetries`<br>`DefaultTimeout` | `MAXRETRIES`<br>`default_timeout` |
| Constant name (Unexported) | camelCase | `maxRetries`<br>`defaultTimeout` | `MAX_RETRIES`<br>`DefaultTimeout` |
| Interface name | Action + -er suffix | `Reader`<br>`Writer` | `ReadInterface`<br>`WriteStuff` |
| Struct name | PascalCase | `UserProfile`<br>`DatabaseConfig` | `userProfile`<br>`database_config` |

## Comment Rules

| Type | Format | Example | Purpose |
| --- | --- | --- | --- |
| Package comment | `// Package [name]` | `// Package handlers provides HTTP handlers.` | Package description |
| Function comment | `// [FuncName]` | `// CreateUser creates a new user` | Exported function documentation |
| Inline comment | `// [Explanation]` | `// retry count` | Code explanation |
| TODO comment | `// TODO:` | `// TODO: implement validation` | Unimplemented / pending items |

## Error Handling

| Rule | Good Example | Bad Example |
| --- | --- | --- |
| Immediate check | `if err != nil {`<br>&nbsp;&nbsp;&nbsp;&nbsp;`return err`<br>`}` | `if err == nil {`<br>&nbsp;&nbsp;&nbsp;&nbsp;`// logic`<br>`}` |
| Contextual error wrapping | `fmt.Errorf("failed to connect: %w", err)` | `return err` |
| Restrict panic usage | Allowed only during initialization | Strictly prohibited in business logic |

## Structs and Methods

| Rule | Good Example | Bad Example |
| --- | --- | --- |
| Field order | 1. Exported fields<br>2. Unexported fields | - |
| Method receiver | Consistent naming<br>(1-2 characters recommended) | Mixing `func (u *User)` and `func (user *User)` |
| Constructor | `New[Type]` format | Mixing `NewUser()` and `CreateUser()` |

## Imports

| Rule | Good Example | Bad Example |
| --- | --- | --- |
| Grouping | 1. Standard library<br>2. Third-party packages<br>3. Internal project packages | Randomly mixed import statements |
| Aliasing | Use only to resolve name collisions | Unnecessary aliasing |

## Concurrency

| Rule | Good Example | Bad Example |
| --- | --- | --- |
| Channel direction | Specify explicitly<br>`chan<- int` (send-only)<br>`<-chan int` (receive-only) | Unspecified direction |
| Goroutine management | Control via `context.Context` | Spawning goroutines without lifecycle control |

## Testing

| Rule | Example | Notes |
| --- | --- | --- |
| File name | `xxx_test.go` | Located in the same package as the code under test |
| Test function name | `Test[FunctionName]` | Clearly describe test scenarios |

## Miscellaneous

| Rule | Description | Example |
| --- | --- | --- |
| Magic numbers | Define as named constants | `const maxRetries = 3` |
| Debugging code | Use `debug` build tag | `//+build debug` |
| Circular dependencies | Strictly forbidden | - |

## Implementation Location Guidelines

| Feature | Implementation Location | Description |
| --- | --- | --- |
| Dropdown / Select APIs | `selection_controller.go` | Endpoints providing selection options for dropdowns should be consolidated in `selection_controller`. |

---

## Detailed Naming Guidelines

Detailed explanations of Go naming conventions:

### 1. Package Names

```go
// Good
package http
package bytes
package strings

// Bad
package HTTP        // Do not use uppercase
package string_util // Do not use underscores
```

- Use only lowercase letters.
- Keep them short and concise.
- Single words are preferred.
- Avoid underscores (`_`) and mixedCaps.
- The package name serves as the qualifier at import sites.

### 2. Interface Names

```go
// Good
type Reader interface { ... }      // Single-method interface
type StringWriter interface { ... } // Descriptive capability name

// Idiomatic -er suffix
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
type Formatter interface { Format(f State, c rune) }
```

### 3. Variable Names

```go
// Local variables: Short and concise
i := 0
sum := 0
done := false

// Global / Package-level variables: More descriptive
var (
	MaxRetryCount  = 3
	DefaultTimeout = 30 * time.Second
)
```

### 4. Constants

```go
// Constants
const (
	StatusOK       = 200
	StatusNotFound = 404

	// Using iota for enumerations
	Sunday = iota
	Monday
	Tuesday
)
```

### 5. Struct Field Names

```go
type User struct {
	ID        int       // Exported, accessible externally
	Name      string
	email     string    // Unexported, internal package usage only
	CreatedAt time.Time
}
```

### 6. Method Names

```go
// Receivers should use short 1-2 character names
func (s *Server) Start() error
func (c *Client) Connect() error

// Method names in CamelCase / PascalCase
func (f *File) ReadLine() (string, error)
func (w *Writer) WriteString(s string) (n int, err error)
```

### 7. Export Rules

```go
// Exported identifiers start with a capital letter
type Customer struct{} // Accessible from outside packages
func ServeHTTP() {}    // Accessible from outside packages

// Unexported identifiers start with a lowercase letter
type customer struct{} // Accessible only within the package
func serveHTTP() {}    // Accessible only within the package
```

### 8. Handling Acronyms and Abbreviations

```go
// Acronyms should be consistently all-uppercase or all-lowercase
ServeHTTP  // Good
ServeHttp  // Bad
XMLParser  // Good
XmlParser  // Bad
```

### 9. Getters and Setters

```go
// Go idiom does not use the "Get" prefix for getters
owner := obj.Owner() // No "Get" prefix
obj.SetOwner(user)   // Use "Set" prefix for setters
```

### Key Principles

1. **Prioritize Readability**:
   - Choose names that clearly convey intent.
   - Use only widely recognized abbreviations.

2. **Length Proportional to Scope**:
   - Short names are suitable for small, local scopes.
   - Use more descriptive names for larger scopes.

3. **Consistency**:
   - Follow consistent naming conventions throughout the project.
   - Adhere to Go standard library conventions.

4. **Consider Context**:
   - Consider how the package name looks at the caller site.
   - Avoid redundant information (e.g., `customer.ID` instead of `customer.CustomerID`).

5. **Formatting**:
   - Always prioritize and enforce code formatting (`go fmt` / `goimports`).

---

## Error Handling

Always use `github.com/cockroachdb/errors` for error handling.  
Its basic usage is identical to the standard library:

```go
errors.New("error message")
```

Handling application errors: Capture/wrap the error, then inspect it using `errors.As`:

```go
if err := con.usecases.LoginUseCase.Execute(); err != nil {
	log.Println(err)
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, gin.H{
			"message": appErr.Message,
		})
		return
	}
	// Fallback for non-application errors
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
	})
	return
}
```

---

## QueryService

The QueryService is strictly responsible for data retrieval (CQRS Read model).  
Cross-table queries across multiple domain entities are permitted here.  
The return value should be `qsdto.Dto` or `dbutil.Paginator`:

```go
type UserQueryService interface {
	FetchUser(ctx context.Context, id uint) (*qsdto.UserDto, error)
	FetchAllUsers(ctx context.Context) ([]qsdto.UserDto, error)
	FetchUsers(ctx context.Context, conditions FetchUsersConditions) (*dbutil.Paginator, error)
}
```

---

## ValueObject

### Criteria for NOT Creating a ValueObject

Do not create ValueObjects for:
- Identifier fields (`ID`)
- Booleans (`bool`)
- Timestamp fields ending with `*At` (`CreatedAt`, `UpdatedAt`), unless specific business logic is involved.

### ValueObject Invariants

ValueObjects must guarantee that they always hold valid values (enforce validation and invariant checks upon creation).
