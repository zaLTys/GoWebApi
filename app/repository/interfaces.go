package repository

import "books-api/app/models"

// BookRepository defines the interface for book data operations
type BookRepository interface {
	Create(book *models.Book) error
	GetByID(id uint) (*models.Book, error)
	GetAll() ([]models.Book, error)
	Update(book *models.Book) error
	Delete(id uint) error
}
