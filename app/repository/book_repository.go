package repository

import (
	"books-api/app/models"
	"gorm.io/gorm"
)

// bookRepository implements the BookRepository interface
type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a new instance of book repository
func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{
		db: db,
	}
}

// Create adds a new book to the database
func (r *bookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

// GetByID retrieves a book by its ID
func (r *bookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// GetAll retrieves all books from the database
func (r *bookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Find(&books).Error
	return books, err
}

// Update modifies an existing book in the database
func (r *bookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

// Delete removes a book from the database by ID
func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}
