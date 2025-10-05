package mocks

import (
	"books-api/app/models"
	"github.com/stretchr/testify/mock"
)

// MockBookService is a mock implementation of BookService interface
type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookService) GetBookByID(id uint) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) GetAllBooks() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookService) UpdateBook(id uint, updateData models.Book) (*models.Book, error) {
	args := m.Called(id, updateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) DeleteBook(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
