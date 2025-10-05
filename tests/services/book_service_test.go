package services_test

import (
	"books-api/app/service"
	"books-api/tests/repositories/mocks"
	"books-api/app/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookService_CreateBook_Success(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	color := models.Red
	book := &models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
		Color:  &color,
	}

	mockRepo.On("Create", book).Return(nil)

	err := svc.CreateBook(book)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookService_CreateBook_InvalidColor(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	invalidColor := models.Color("Purple")
	book := &models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
		Color:  &invalidColor,
	}

	err := svc.CreateBook(book)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid color")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestBookService_GetBookByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	expectedBook := &models.Book{
		ID:     1,
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
	}

	mockRepo.On("GetByID", uint(1)).Return(expectedBook, nil)

	book, err := svc.GetBookByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedBook, book)
	mockRepo.AssertExpectations(t)
}

func TestBookService_GetBookByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("record not found"))

	book, err := svc.GetBookByID(999)
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "book not found")
	mockRepo.AssertExpectations(t)
}

func TestBookService_GetAllBooks(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	expectedBooks := []models.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1", Pages: 100},
		{ID: 2, Title: "Book 2", Author: "Author 2", Pages: 200},
	}

	mockRepo.On("GetAll").Return(expectedBooks, nil)

	books, err := svc.GetAllBooks()
	assert.NoError(t, err)
	assert.Len(t, books, 2)
	mockRepo.AssertExpectations(t)
}

func TestBookService_UpdateBook_Success(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:     1,
		Title:  "Original Title",
		Author: "Original Author",
		Pages:  100,
	}

	updateData := models.Book{
		Title: "Updated Title",
		Pages: 150,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Book")).Return(nil)

	book, err := svc.UpdateBook(1, updateData)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", book.Title)
	assert.Equal(t, 150, book.Pages)
	mockRepo.AssertExpectations(t)
}

func TestBookService_UpdateBook_InvalidColor(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:     1,
		Title:  "Original Title",
		Author: "Original Author",
		Pages:  100,
	}

	invalidColor := models.Color("Yellow")
	updateData := models.Book{
		Color: &invalidColor,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil)

	book, err := svc.UpdateBook(1, updateData)
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "invalid color")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestBookService_DeleteBook_Success(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	existingBook := &models.Book{
		ID:     1,
		Title:  "To Be Deleted",
		Author: "Test Author",
		Pages:  100,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingBook, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.DeleteBook(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBookService_DeleteBook_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	svc := service.NewBookService(mockRepo)

	mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("record not found"))

	err := svc.DeleteBook(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "book not found")
	mockRepo.AssertNotCalled(t, "Delete")
}
