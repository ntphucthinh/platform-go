# OpenAPI Design Guidelines

## Introduction

This document outlines the conventions and rules for API design in this project, including OpenAPI definition methods, naming conventions, and directory structures. API designers and developers must adhere to this document when defining OpenAPI specifications.

## 1. Basic File Structure

- **Root file**: `src/openapi/openapi.yml`
  - The entry-point root file that aggregates all OpenAPI definitions.
- **Path definitions**: `src/openapi/paths/`
  - Contains API endpoint definitions (HTTP methods such as GET, POST, PUT, DELETE, etc.).
- **Component definitions**: `src/openapi/components/`
  - Defines reusable schemas (data models), parameters, responses, etc.
    - `schemas/`: Defines data structures for request bodies and responses.

## 2. Naming Conventions

To ensure high readability and maintainability of APIs, follow these naming conventions:

### 2.1. File Names

File names related to OpenAPI definitions use **lowercase snake_case** (`lower_snake_case`).

- Examples: `user_profile.yml`, `error_response.yml`

### 2.2. Component Names

Schema names, parameter names, response names, etc. defined under the `components` section use **UpperCamelCase** (PascalCase).

- Examples: `UserProfile`, `ErrorResponse`, `UserIDParameter`

### 2.3. operationId

The `operationId` uniquely identifying each API operation must follow this format:

- **Format**: `Verb` + `ResourceName` (e.g., `FetchUsers`, `CreateUser`)
- **Verb Examples**:
  - Retrieval (single / multiple): `Fetch` (e.g., `FetchUserByID`, `FetchUsers`)
  - Creation: `Create` (e.g., `CreateUser`)
  - Update: `Update` (e.g., `UpdateUser`)
  - Deletion: `Delete` (e.g., `DeleteUser`)
- **Special Path Naming Conventions**:
  - For `/selections/<route_name>`: `Selection<RouteName>` (e.g., the `operationId` for `/selections/Prefectures` is `SelectionPrefectures`)

### 2.4. Parameter Names in `parameters`

The `name` field of parameters referenced in `parameters` must be written in **camelCase**.

- **Examples**:
  - `userID`
  - `pageSize`
  - `isComplete`

### 2.5. `components/schemas` References

When referencing `requestBody` or `responses` schemas inside `paths` definition files, use the following naming convention:

- **Reference format**: `#/components/schemas/{operationId}{Request|Response}`
- **Examples**:
  - When `operationId` is `FetchUserByID`:
    - Request Body: `./components/schemas/users/:id.yml#/FetchUserByIDRequest`
    - Response: `./components/schemas/users/:id.yml#/FetchUserByIDResponse`

This ensures that each operation and its corresponding request/response schemas are immediately associated with each other.

**Exception**: When clearer contextual names are required (such as authentication-related endpoints), deviations are permitted (e.g., `AuthRequest`, `AuthResponse`). In general, adhere to the standard naming rule above.

### 2.6. `parameters` References

When referencing `parameters` inside `paths` definition files, use the following naming convention:

- **Reference format**: `./components/schemas/{resource_name}/{parameter_file_name}.yml#/{operationId}-{ParameterName}`
  - **Note**: The actual file path must comply with the rules outlined in "3. Directory Structure" below. This format represents the logical structure.
- **Examples**:
  - When `operationId` is `FetchUsers` and parameter name is `Name` (defined in the `users` resource `index.yml`):
    - Parameter reference: `$ref: ../components/schemas/users/index.yml#/FetchUsers-Name`
  - When `operationId` is `FetchUserByID` and parameter name is `UserID` (defined in the `users` resource `:id.yml`):
    - Parameter reference: `$ref: ../components/schemas/users/:id.yml#/FetchUserByID-UserID`

### 2.7. Property Names

Property names defined in `components/schemas` must use **lowerCamelCase**.

- Examples: `userName`, `projectStatus`, `isArchived`

## 3. Directory Structure

OpenAPI definition files are organized to reflect the API routing structure. This allows intuitive discovery of the definition file corresponding to any specific API endpoint. This rule applies to both the `paths` directory (endpoint definitions) and subdirectories inside `components` (such as `schemas`).

### 3.1. `paths` Directory

Create directories and files mirroring the API URL path structure:

- **Root path or single-level path** (e.g., `/users`):
  - Create a directory named after the path and place `index.yml` directly inside it.
  - Example: `paths/users/index.yml` (definition for the `/users` endpoint)
- **Paths containing path parameters** (e.g., `/users/{userID}`):
  - Represent parameter segments using placeholder notation (e.g., `:userID`) as the file or directory name.
  - Example: `paths/users/:userID.yml` (definition for the `/users/{userID}` endpoint)
- **Nested paths** (e.g., `/projects/{projectID}/tasks`):
  - Create a directory for each segment, naming the final file after the leaf segment.
  - Example: `paths/projects/:projectID/tasks.yml` (definition for the `/projects/{projectID}/tasks` endpoint)
- **Exceptions**:
  - `paths/auth.yml` (`/auth` related endpoints)
  - `paths/health.yml` (`/health` endpoint)

### 3.2. `components/schemas` Directory

Follows the same structural rules as the `paths` directory. In addition, the following special directories exist:

- `components/schemas/common/`: Defines schemas shared across multiple APIs (e.g., `error_response.yml`).
- `components/schemas/enum/`: Defines enum types (e.g., `user_status.yml`).  
  **When using enums, always define them inside this directory.**

## 4. Data Type Definition Rules

### 4.1. Handling Date and Time Data

Adhere to the following rules when defining date and time fields:

#### 4.1.1. Date and Time in Responses

Response date/time fields should use `string` rather than `date` or `date-time` types:

- Date (`date`)
  - Format: `2006-01-02`
  - Example: `time.Now().Format(time.DateOnly)`
- Date-Time (`date-time`)
  - Format: `2006-01-02 15:04:05`
  - Example: `time.Now().Format(time.DateTime)`

```yaml
# Response schema example
createdAt:
  type: string
  description: "Creation timestamp"
  example: "2023-01-15 10:30:00"

updatedAt:
  type: string
  description: "Update timestamp"
  example: "2023-01-15 14:45:00"
```

#### 4.1.2. Date and Time in Requests

Request datetime fields should use custom Go types instead of standard `date-time`:

```yaml
# Request schema example
startDate:
  type: string
  x-go-type: types.CustomTime
  description: "Start timestamp"
  example: "2023-01-15 10:30:00"

endDate:
  type: string
  x-go-type: types.CustomTime
  description: "End timestamp"
  example: "2023-01-15 18:00:00"
```

**Notes**:

- Before using `x-go-type: types.CustomTime` in requests, verify that `types.CustomTime` is properly implemented.
- Always use the uniform format `2006-01-02 15:04:05` (Go standard time layout).
