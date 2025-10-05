package controllers_test

import (
	"books-api/app/controller"
	"books-api/tests/services/mocks"
	"books-api/app/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestBookController_CreateBook_Success(t *testing.T) {
	mockService := new(mocks.MockBookService)
	ctrl := controller.NewBookController(mockService)
	router := setupTestRouter()

	router.POST("/books", ctrl.CreateBook)

	color := models.Red
	book := models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
		Color:  &color,
	}

	mockService.On("CreateBook", &book).Return(nil)

	body, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookController_CreateBook_InvalidJSON(t *testing.T) {
	mockService := new(mocks.MockBookService)
	ctrl := controller.NewBookController(mockService)
	router := setupTestRouter()

	router.POST("/books", ctrl.CreateBook)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBookController_CreateBook_ValidationError(t *testing.T) {
	mockService := new(mocks.MockBookService)
	ctrl := controller.NewBookController(mockService)
	router := setupTestRouter()

	router.POST("/books", ctrl.CreateBook)

	invalidColor := models.Color("Purple")
	book := models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
		Color:  &invalidColor,
	}

	mockService.On("CreateBook", &book).Return(errors.New("invalid color: Purple"))

	body, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookController_ListBooks(t *testing.T) {
	mockService := new(mocks.MockBookService)
	ctrl := controller.NewBookController(mockService)
	router := setupTestRouter()

	router.GET("/books", ctrl.ListBooks)

	expectedBooks := []models.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1", Pages: 100},
		{ID: 2, Title: "Book 2", Author: "Author 2", Pages: 200},
	}

	mockService.On("GetAllBooks").Return(expectedBooks, nil)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var books []models.Book
	json.Unmarshal(w.Body.Bytes(), &books)
	assert.Len(t, books, 2)
	mockService.AssertExpectations(t)
}

func TestBookController_GetBook_Success(t *testing.T) {
	mockService := new(mocks.MockBookService)
	ctrl := controller.NewBookController(mockService)
	router := setupTestRouter()

	router.GET("/books/:id", ctrl.GetBook)

	expectedBook := &models.Book{
		ID:     1,
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
	}

	mockService.On("GetBookByID", uint(1)).Return(expectedBook, nil)

	req, _ := http.NewRequest("GET", "/books/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var book models.Book
	json.Unmarshal(w.Body.Bytes(), &book)
	assert.Equal(t, "Test Book", book.Title)
	mockService.AssertExpectations(t)
}

func TestBookController_GetBook_NotFound(t *testing.T) {
	mockService := new(mocks.MockBookService)
	ctrl := controller.NewBookController(mockService)
	router := setupTestRouter()

	router.GET("/books/:id", ctrl.GetBook)

	mockService.On("GetBookByID", uint(999)).Return(nil, errors.New("book not found"))

	req, _ := http.NewRequest("GET", "/books/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookController_UpdateBook_Success(t *testing.T) {
	mockService := new(mocks.MockBookService)
	ctrl := controller.NewBookController(mockService)
	router := setupTestRouter()

	router.PUT("/books/:id", ctrl.UpdateBook)

	updateData := models.Book{
		Title: "Updated Title",
		Pages: 150,
	}

	updatedBook := &models.Book{
		ID:     1,
		Title:  "Updated Title",
		Author: "Original Author",
		Pages:  150,
	}

	mockService.On("UpdateBook", uint(1), updateData).Return(updatedBook, nil)

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/books/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBookController_DeleteBook_Success(t *testing.T) {
	mockService := new(mocks.MockBookService)
	ctrl := controller.NewBookController(mockService)
	router := setupTestRouter()

	router.DELETE("/books/:id", ctrl.DeleteBook)

	mockService.On("DeleteBook", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/books/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
