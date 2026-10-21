# Testing Scope

## Mandatory Test Targets

- **UseCase**
  - UseCase
    - `src/internal/usecases/***/usecase`
- **Domain**
  - ValueObject
    - `src/internal/domain/***/valueobjects`
  - Entity
    - `src/internal/domain/***/entity`
  - Service
    - `src/internal/domain/***/service`

## Coverage

The mandatory test targets listed above must exceed **60%** code coverage.  
This 60% requirement does not apply to overall project-wide coverage.

## Urgent Implementations

If rapid delivery or hotfix implementation is required, ensure at least **one** happy-path (normal) test case is created.  
Having at least one test case makes it significantly easier for other team members to extend and add more tests later.
