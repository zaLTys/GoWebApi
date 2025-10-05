package service

import (
	"books-api/app/repository"
	"books-api/app/models"
	"fmt"
	"log"
)

// bookService implements the BookService interface
type bookService struct {
	bookRepo repository.BookRepository
}

// NewBookService creates a new instance of book service
func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

// CreateBook creates a new book with validation and logging
func (s *bookService) CreateBook(book *models.Book) error {
	log.Printf("Creating new book: %s by %s", book.Title, book.Author)
	
	// Validate color if provided
	if book.Color != nil && !book.Color.IsValid() {
		log.Printf("Invalid color provided for book: %s", *book.Color)
		return fmt.Errorf("invalid color: %s", *book.Color)
	}
	
	err := s.bookRepo.Create(book)
	if err != nil {
		log.Printf("Failed to create book: %v", err)
		return fmt.Errorf("failed to create book: %w", err)
	}
	
	log.Printf("Successfully created book with ID: %d", book.ID)
	return nil
}

// GetBookByID retrieves a book by ID with logging
func (s *bookService) GetBookByID(id uint) (*models.Book, error) {
	log.Printf("Retrieving book with ID: %d", id)
	
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		log.Printf("Failed to retrieve book with ID %d: %v", id, err)
		return nil, fmt.Errorf("book not found")
	}
	
	log.Printf("Successfully retrieved book: %s", book.Title)
	return book, nil
}

// GetAllBooks retrieves all books with logging
func (s *bookService) GetAllBooks() ([]models.Book, error) {
	log.Printf("Retrieving all books")
	
	books, err := s.bookRepo.GetAll()
	if err != nil {
		log.Printf("Failed to retrieve books: %v", err)
		return nil, fmt.Errorf("failed to retrieve books: %w", err)
	}
	
	log.Printf("Successfully retrieved %d books", len(books))
	return books, nil
}

// UpdateBook updates an existing book with validation and logging
func (s *bookService) UpdateBook(id uint, updateData models.Book) (*models.Book, error) {
	log.Printf("Updating book with ID: %d", id)
	
	// First, get the existing book
	existingBook, err := s.bookRepo.GetByID(id)
	if err != nil {
		log.Printf("Book with ID %d not found for update: %v", id, err)
		return nil, fmt.Errorf("book not found")
	}
	
	// Validate color if provided in update
	if updateData.Color != nil && !updateData.Color.IsValid() {
		log.Printf("Invalid color provided for book update: %s", *updateData.Color)
		return nil, fmt.Errorf("invalid color: %s", *updateData.Color)
	}
	
	// Update the existing book with new data
	existingBook.Update(updateData)
	
	err = s.bookRepo.Update(existingBook)
	if err != nil {
		log.Printf("Failed to update book with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to update book: %w", err)
	}
	
	log.Printf("Successfully updated book: %s", existingBook.Title)
	return existingBook, nil
}

// DeleteBook deletes a book by ID with logging
func (s *bookService) DeleteBook(id uint) error {
	log.Printf("Deleting book with ID: %d", id)
	
	// Check if book exists first
	_, err := s.bookRepo.GetByID(id)
	if err != nil {
		log.Printf("Book with ID %d not found for deletion: %v", id, err)
		return fmt.Errorf("book not found")
	}
	
	err = s.bookRepo.Delete(id)
	if err != nil {
		log.Printf("Failed to delete book with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete book: %w", err)
	}
	
	log.Printf("Successfully deleted book with ID: %d", id)
	return nil
}
