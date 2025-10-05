# Books API - Layered Architecture

This project demonstrates a clean layered architecture for a Go web API using Gin, GORM, and Swagger.

## Architecture Overview

The application follows a layered architecture pattern with clear separation of concerns:

```
┌─────────────────┐
│   Controllers   │  ← HTTP handlers, request/response mapping
├─────────────────┤
│    Services     │  ← Business logic, validation, logging
├─────────────────┤
│  Repositories   │  ← Data access layer, database operations
├─────────────────┤
│    Database     │  ← SQLite database
└─────────────────┘
```

## Project Structure

```
.
├── main.go                           # Application entry point
├── models/                           # Domain models
│   ├── models.go                     # Book and Color models
│   └── models_test.go                # Model unit tests
├── app/                              # Application code
│   ├── controller/                   # HTTP handlers
│   │   └── book_controller.go        # Book API endpoints
│   ├── service/                      # Business logic layer
│   │   ├── interfaces.go             # Service interfaces
│   │   └── book_service.go           # Book business logic with logging
│   ├── repository/                   # Data access layer
│   │   ├── interfaces.go             # Repository interfaces
│   │   └── book_repository.go        # Book database operations
│   └── migrations/                   # Database migrations
│       ├── interfaces.go             # Migration interfaces
│       └── migration_manager.go      # GORM migration management
├── tests/                            # All test files
│   ├── app/                          # Unit tests (mirrors app structure)
│   │   ├── controller/               # Controller tests
│   │   │   └── book_controller_test.go
│   │   ├── service/                  # Service tests
│   │   │   ├── book_service_test.go
│   │   │   └── mocks/                # Service mocks
│   │   │       └── book_service_mock.go
│   │   └── repository/               # Repository tests
│   │       ├── book_repository_test.go
│   │       └── mocks/                # Repository mocks
│   │           └── book_repository_mock.go
│   └── integration/                  # Integration tests
│       └── book_api_test.go          # Full stack E2E tests
├── docs/                             # Swagger documentation
│   ├── docs.go                       # Generated Swagger docs
│   ├── swagger.json                  # Swagger JSON spec
│   └── swagger.yaml                  # Swagger YAML spec
├── go.mod                            # Go module definition
├── Makefile                          # Build and test commands
└── README.md                         # Project documentation
```

## Layer Responsibilities

### 1. Controllers (`app/controller/`)
- Handle HTTP requests and responses
- Input validation and parameter parsing
- Route request to appropriate service methods
- Return appropriate HTTP status codes

### 2. Services (`app/service/`)
- Implement business logic and validation rules
- Log all operations for audit and debugging
- Coordinate between different repositories if needed
- Handle error translation and formatting

### 3. Repositories (`app/repository/`)
- Abstract database operations behind interfaces
- Implement CRUD operations using GORM
- Handle database-specific error handling
- Provide clean data access API

### 4. Migrations (`app/migrations/`)
- Manage database schema changes
- Provide interface for running migrations
- Centralize database initialization logic

## Key Features

### Interface-Based Design
All layers are defined by interfaces, making the code:
- Testable (easy to mock dependencies)
- Flexible (easy to swap implementations)
- Maintainable (clear contracts between layers)

### Comprehensive Logging
The service layer logs all operations including:
- Create, read, update, delete operations
- Validation errors
- Success and failure outcomes
- Performance metrics

### Clean Error Handling
- Repository layer returns raw errors
- Service layer translates to business errors
- Controller layer maps to appropriate HTTP responses

## API Endpoints

| Method | Endpoint      | Description           |
|--------|---------------|-----------------------|
| GET    | /books        | List all books        |
| POST   | /books        | Create a new book     |
| GET    | /books/{id}   | Get book by ID        |
| PUT    | /books/{id}   | Update book by ID     |
| DELETE | /books/{id}   | Delete book by ID     |
| GET    | /swagger/*    | Swagger documentation |

## Running the Application

1. **Start the server:**
   ```bash
   go run .
   ```

2. **Access the API:**
   - API Base URL: `http://localhost:8080`
   - Swagger UI: `http://localhost:8080/swagger/index.html`

3. **Example API calls:**
   ```bash
   # Get all books
   curl http://localhost:8080/books
   
   # Create a new book
   curl -X POST http://localhost:8080/books \
     -H "Content-Type: application/json" \
     -d '{"title":"Foundation","author":"Isaac Asimov","pages":255,"color":"Blue"}'
   
   # Get book by ID
   curl http://localhost:8080/books/1
   ```

## Development Benefits

### 1. **Testability**
Each layer can be unit tested independently using interface mocks.

### 2. **Maintainability**
Clear separation of concerns makes code easy to understand and modify.

### 3. **Scalability**
New features can be added by extending interfaces and implementing new services.

### 4. **Flexibility**
Database or framework changes only affect specific layers.

### 5. **Observability**
Comprehensive logging in the service layer provides excellent debugging capabilities.

## Future Enhancements

- Add middleware for authentication and authorization
- Implement caching in the service layer
- Add metrics and monitoring
- Implement database connection pooling
- Add integration tests for the full stack
