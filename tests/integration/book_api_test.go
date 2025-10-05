package integration

import (
	"books-api/app/controller"
	"books-api/app/migrations"
	"books-api/app/repository"
	"books-api/app/service"
	"books-api/app/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BookAPITestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
}

func (suite *BookAPITestSuite) SetupTest() {
	// Create in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.NoError(err)

	// Run migrations
	migrationManager := migrations.NewMigrationManager()
	err = migrationManager.RunMigrations(db)
	suite.NoError(err)

	// Setup layers
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookController := controller.NewBookController(bookService)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	bookRoutes := router.Group("/books")
	{
		bookRoutes.POST("", bookController.CreateBook)
		bookRoutes.GET("", bookController.ListBooks)
		bookRoutes.GET("/:id", bookController.GetBook)
		bookRoutes.PUT("/:id", bookController.UpdateBook)
		bookRoutes.DELETE("/:id", bookController.DeleteBook)
	}

	suite.db = db
	suite.router = router
}

func (suite *BookAPITestSuite) TearDownTest() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func (suite *BookAPITestSuite) TestCreateBook_Success() {
	color := models.Red
	book := map[string]interface{}{
		"title":  "Integration Test Book",
		"author": "Test Author",
		"pages":  250,
		"color":  color,
	}

	body, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response models.Book
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Integration Test Book", response.Title)
	assert.NotZero(suite.T(), response.ID)
}

func (suite *BookAPITestSuite) TestCreateBook_InvalidColor() {
	book := map[string]interface{}{
		"title":  "Invalid Book",
		"author": "Test Author",
		"pages":  100,
		"color":  "Purple",
	}

	body, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *BookAPITestSuite) TestGetAllBooks() {
	// Create some books first
	books := []models.Book{
		{Title: "Book 1", Author: "Author 1", Pages: 100},
		{Title: "Book 2", Author: "Author 2", Pages: 200},
		{Title: "Book 3", Author: "Author 3", Pages: 300},
	}

	for _, book := range books {
		suite.db.Create(&book)
	}

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response []models.Book
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(suite.T(), response, 3)
}

func (suite *BookAPITestSuite) TestGetBookByID_Success() {
	// Create a book
	book := models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  150,
	}
	suite.db.Create(&book)

	req, _ := http.NewRequest("GET", "/books/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.Book
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Test Book", response.Title)
}

func (suite *BookAPITestSuite) TestGetBookByID_NotFound() {
	req, _ := http.NewRequest("GET", "/books/999", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *BookAPITestSuite) TestUpdateBook_Success() {
	// Create a book
	book := models.Book{
		Title:  "Original Title",
		Author: "Original Author",
		Pages:  100,
	}
	suite.db.Create(&book)

	updateData := map[string]interface{}{
		"title": "Updated Title",
		"pages": 200,
	}

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/books/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.Book
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "Updated Title", response.Title)
	assert.Equal(suite.T(), 200, response.Pages)
}

func (suite *BookAPITestSuite) TestDeleteBook_Success() {
	// Create a book
	book := models.Book{
		Title:  "To Be Deleted",
		Author: "Test Author",
		Pages:  100,
	}
	suite.db.Create(&book)

	req, _ := http.NewRequest("DELETE", "/books/1", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Verify deletion
	var count int64
	suite.db.Model(&models.Book{}).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *BookAPITestSuite) TestCompleteWorkflow() {
	// 1. Create a book
	color := models.Green
	createData := map[string]interface{}{
		"title":  "Workflow Book",
		"author": "Workflow Author",
		"pages":  300,
		"color":  color,
	}

	body, _ := json.Marshal(createData)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createdBook models.Book
	json.Unmarshal(w.Body.Bytes(), &createdBook)

	// 2. Get the book
	req, _ = http.NewRequest("GET", "/books/1", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 3. Update the book
	updateData := map[string]interface{}{
		"pages": 350,
	}
	body, _ = json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", "/books/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 4. Delete the book
	req, _ = http.NewRequest("DELETE", "/books/1", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 5. Verify it's gone
	req, _ = http.NewRequest("GET", "/books/1", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func TestBookAPITestSuite(t *testing.T) {
	suite.Run(t, new(BookAPITestSuite))
}
