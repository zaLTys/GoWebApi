package service

import "books-api/app/models"

// BookService defines the interface for book business logic
type BookService interface {
	CreateBook(book *models.Book) error
	GetBookByID(id uint) (*models.Book, error)
	GetAllBooks() ([]models.Book, error)
	UpdateBook(id uint, updateData models.Book) (*models.Book, error)
	DeleteBook(id uint) error
}
