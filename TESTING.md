# Testing Infrastructure

This project includes comprehensive testing at all layers of the architecture.

## Test Coverage Summary

| Layer | Coverage | Test Type |
|-------|----------|-----------|
| **Repository** | 100% | Unit Tests (in-memory DB) |
| **Service** | 80.4% | Unit Tests (mocked repository) |
| **Controller** | 64.7% | Unit Tests (mocked service) |
| **Models** | 95.5% | Unit Tests |
| **Integration** | Full E2E | Integration Tests |

## Running Tests

### All Tests
```bash
make test
```

### Unit Tests Only
```bash
make test-unit
```

### Integration Tests Only
```bash
make test-integration
```

### With Coverage Report
```bash
make test-coverage
# Opens coverage.html in browser
```

### With Race Detector
```bash
make test-race
```

## Test Structure

```
.
├── app/                                 # Application code (no tests)
│   ├── controller/
│   │   └── book_controller.go
│   ├── repository/
│   │   └── book_repository.go
│   └── service/
│       └── book_service.go
├── models/
│   ├── models.go
│   └── models_test.go                   # Model unit tests
└── tests/                               # All tests organized here
    ├── app/                             # Unit tests (mirrors app structure)
    │   ├── controller/
    │   │   └── book_controller_test.go  # Controller unit tests
    │   ├── repository/
    │   │   ├── book_repository_test.go  # Repository unit tests
    │   │   └── mocks/
    │   │       └── book_repository_mock.go
    │   └── service/
    │       ├── book_service_test.go     # Service unit tests
    │       └── mocks/
    │           └── book_service_mock.go
    └── integration/
        └── book_api_test.go             # Full stack integration tests
```

## Testing Patterns

### 1. Repository Tests (Database Layer)
**Pattern**: Use in-memory SQLite database

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(t, err)
    err = db.AutoMigrate(&models.Book{})
    assert.NoError(t, err)
    return db
}
```

**What we test**:
- ✅ CRUD operations
- ✅ Error handling (not found, etc.)
- ✅ Database constraints
- ✅ Query correctness

### 2. Service Tests (Business Logic Layer)
**Pattern**: Use mocked repository

```go
mockRepo := new(mocks.MockBookRepository)
service := NewBookService(mockRepo)

mockRepo.On("Create", book).Return(nil)
err := service.CreateBook(book)
mockRepo.AssertExpectations(t)
```

**What we test**:
- ✅ Business logic validation
- ✅ Error translation
- ✅ Logging behavior
- ✅ Edge cases

### 3. Controller Tests (HTTP Layer)
**Pattern**: Use mocked service + httptest

```go
mockService := new(mocks.MockBookService)
controller := NewBookController(mockService)
router := setupTestRouter()

mockService.On("CreateBook", &book).Return(nil)
req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
```

**What we test**:
- ✅ HTTP status codes
- ✅ Request/response JSON
- ✅ Input validation
- ✅ Error responses

### 4. Integration Tests (Full Stack)
**Pattern**: Test suite with real database and full stack

```go
type BookAPITestSuite struct {
    suite.Suite
    db     *gorm.DB
    router *gin.Engine
}

func (suite *BookAPITestSuite) SetupTest() {
    // Create fresh DB and full stack for each test
}
```

**What we test**:
- ✅ Complete user workflows
- ✅ All layers working together
- ✅ Database persistence
- ✅ Real HTTP requests

## Test Examples

### Unit Test Example
```go
func TestBookService_CreateBook_InvalidColor(t *testing.T) {
    mockRepo := new(mocks.MockBookRepository)
    service := NewBookService(mockRepo)

    invalidColor := models.Color("Purple")
    book := &models.Book{
        Title:  "Test Book",
        Color:  &invalidColor,
    }

    err := service.CreateBook(book)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid color")
    mockRepo.AssertNotCalled(t, "Create")
}
```

### Integration Test Example
```go
func (suite *BookAPITestSuite) TestCompleteWorkflow() {
    // 1. Create
    // 2. Read
    // 3. Update
    // 4. Delete
    // 5. Verify deletion
}
```

## Mocking Strategy

We use `testify/mock` for creating mocks:

**Benefits**:
- ✅ Type-safe mocking
- ✅ Expectation verification
- ✅ Easy to maintain
- ✅ Clear test failures

**Mock Interfaces**:
- `MockBookRepository` - Mocks database operations
- `MockBookService` - Mocks business logic

## Best Practices

### 1. Test Isolation
- Each test should be independent
- Use `SetupTest` and `TearDownTest` for test suites
- Clean database state between tests

### 2. Test Naming
```go
func TestComponent_Method_Scenario(t *testing.T)
// Example: TestBookService_CreateBook_InvalidColor
```

### 3. Arrange-Act-Assert Pattern
```go
func TestExample(t *testing.T) {
    // Arrange: Setup test data and mocks
    mockRepo := new(mocks.MockBookRepository)
    
    // Act: Execute the function under test
    result, err := service.DoSomething()
    
    // Assert: Verify the results
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### 4. Table-Driven Tests
```go
tests := []struct {
    name  string
    input string
    want  bool
}{
    {"valid case", "Red", true},
    {"invalid case", "Purple", false},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := IsValid(tt.input)
        assert.Equal(t, tt.want, got)
    })
}
```

## Continuous Integration

Add to your CI/CD pipeline:

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      - run: make test
      - run: make test-race
```

## Test Metrics

Current test statistics:
- **Total Tests**: 32
- **Unit Tests**: 24
- **Integration Tests**: 8
- **Average Execution Time**: ~2 seconds
- **Repository Coverage**: 100%
- **Service Coverage**: 80.4%
- **Controller Coverage**: 64.7%
- **Models Coverage**: 95.5%

## Future Enhancements

- [ ] Add benchmark tests for performance-critical paths
- [ ] Add mutation testing
- [ ] Add contract testing for API
- [ ] Add load/stress tests
- [ ] Add E2E tests with real database
- [ ] Add snapshot testing for API responses
