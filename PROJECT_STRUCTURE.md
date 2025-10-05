# Project Structure

## Overview

This project follows a clean, organized structure with clear separation between application code and tests.

## Directory Layout

```
GoWebApi/
├── app/                              # Application code (production)
│   ├── controller/                   # HTTP request handlers
│   ├── service/                      # Business logic layer
│   ├── repository/                   # Data access layer
│   ├── models/                       # Domain models
│   │   └── models.go                 # Book and Color models
│   └── migrations/                   # Database migrations
│
├── tests/                            # All test files
│   ├── controllers/                  # Controller unit tests
│   ├── services/                     # Service unit tests + mocks
│   ├── repositories/                 # Repository unit tests + mocks
│   ├── models/                       # Model unit tests
│   │   └── models_test.go            # Model tests
│   └── integration/                  # End-to-end integration tests
│
├── docs/                             # Auto-generated Swagger docs
├── main.go                           # Application entry point
├── Makefile                          # Build and test commands
├── go.mod                            # Go module definition
└── README.md                         # Project documentation
```

## Design Principles

### 1. **Separation of Concerns**
- **app/**: Contains only production code
- **tests/**: Contains all test files and mocks
- **models/**: Domain models with their tests

### 2. **Flat Test Structure**
Tests are organized by layer type with plural names:
```
app/controller/book_controller.go
tests/controllers/book_controller_test.go

app/service/book_service.go
tests/services/book_service_test.go

app/repository/book_repository.go
tests/repositories/book_repository_test.go
```

This makes it easy to find tests and keeps the structure clean.

### 3. **Test Organization**
- **Unit Tests**: In `tests/{controllers,services,repositories,models}/` - test individual components
- **Integration Tests**: In `tests/integration/` - test full stack workflows
- **Mocks**: Stored with their respective test layers

### 4. **Mock Management**
Mocks are stored with their respective test layers:
```
tests/services/mocks/book_service_mock.go
tests/repositories/mocks/book_repository_mock.go
```

## Benefits of This Structure

### ✅ **Clean Production Code**
- No test files mixed with production code
- Clear what's shipped vs what's for testing
- Easier to package and deploy

### ✅ **Easy Navigation**
- Tests organized by layer type
- Find tests quickly: `app/controller/file.go` → `tests/controllers/file_test.go`
- Logical grouping of related code

### ✅ **Better IDE Support**
- GoLand/VSCode can easily distinguish test vs production
- Test coverage tools work better
- Easier to exclude tests from builds

### ✅ **Scalability**
- Easy to add new features (add to `app/` and `tests/app/`)
- Clear where new tests should go
- Mocks are organized and findable

## Package Naming Convention

### Production Code
```go
package controller  // app/controller/
package service     // app/service/
package repository  // app/repository/
package models      // app/models/
```

### Test Code
```go
package controllers_test  // tests/controllers/
package services_test     // tests/services/
package repositories_test // tests/repositories/
package models_test       // tests/models/
```

Using `_test` suffix allows tests to:
- Import the package they're testing
- Test public API only (black-box testing)
- Avoid circular dependencies

## Running Tests

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests only
make test-integration

# With coverage
make test-coverage
```

## Adding New Features

### 1. Add Production Code
```bash
# Create new feature in app/
app/feature/
  ├── feature.go
  └── interfaces.go
```

### 2. Add Tests
```bash
# Create corresponding tests
tests/features/
  ├── feature_test.go
  └── mocks/
      └── feature_mock.go
```

### 3. Add Integration Tests (if needed)
```bash
tests/integration/
  └── feature_integration_test.go
```

## Comparison with Other Structures

### vs. Tests Next to Code
```
❌ Old:
app/controller/
  ├── book_controller.go
  └── book_controller_test.go

✅ New:
app/controller/
  └── book_controller.go
tests/app/controller/
  └── book_controller_test.go
```

**Benefits**:
- Cleaner production directories
- Better separation of concerns
- Easier to exclude tests from builds

### vs. Flat Test Directory
```
❌ Flat:
tests/
  ├── book_controller_test.go
  ├── book_service_test.go
  └── book_repository_test.go

✅ Structured:
tests/app/
  ├── controller/book_controller_test.go
  ├── service/book_service_test.go
  └── repository/book_repository_test.go
```

**Benefits**:
- Mirrors production structure
- Scales better with many files
- Easier to navigate

## Best Practices

1. **Keep tests close conceptually**: Tests mirror production structure
2. **One test file per production file**: Easy 1:1 mapping
3. **Mocks with tests**: Store mocks near the tests that use them
4. **Integration tests separate**: Different concerns, different directory
5. **Model tests with models**: Domain logic stays with domain

## Summary

This structure provides:
- ✅ Clear separation between production and test code
- ✅ Easy navigation and discoverability
- ✅ Scalable organization
- ✅ Better tooling support
- ✅ Professional project layout

It's a modern, maintainable structure suitable for production Go applications!
