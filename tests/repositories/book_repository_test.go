package repositories_test

import (
	"books-api/app/repository"
	"books-api/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&models.Book{})
	assert.NoError(t, err)

	return db
}

func TestBookRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookRepository(db)

	color := models.Red
	book := &models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  100,
		Color:  &color,
	}

	err := repo.Create(book)
	assert.NoError(t, err)
	assert.NotZero(t, book.ID)
}

func TestBookRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookRepository(db)

	// Create a book first
	color := models.Blue
	book := &models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Pages:  200,
		Color:  &color,
	}
	err := repo.Create(book)
	assert.NoError(t, err)

	// Retrieve it
	retrieved, err := repo.GetByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, book.Title, retrieved.Title)
	assert.Equal(t, book.Author, retrieved.Author)
	assert.Equal(t, book.Pages, retrieved.Pages)
}

func TestBookRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookRepository(db)

	_, err := repo.GetByID(999)
	assert.Error(t, err)
}

func TestBookRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookRepository(db)

	// Create multiple books
	books := []*models.Book{
		{Title: "Book 1", Author: "Author 1", Pages: 100},
		{Title: "Book 2", Author: "Author 2", Pages: 200},
		{Title: "Book 3", Author: "Author 3", Pages: 300},
	}

	for _, book := range books {
		err := repo.Create(book)
		assert.NoError(t, err)
	}

	// Retrieve all
	allBooks, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, allBooks, 3)
}

func TestBookRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookRepository(db)

	// Create a book
	color := models.Green
	book := &models.Book{
		Title:  "Original Title",
		Author: "Original Author",
		Pages:  100,
		Color:  &color,
	}
	err := repo.Create(book)
	assert.NoError(t, err)

	// Update it
	book.Title = "Updated Title"
	book.Pages = 150
	err = repo.Update(book)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, 150, updated.Pages)
}

func TestBookRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBookRepository(db)

	// Create a book
	book := &models.Book{
		Title:  "To Be Deleted",
		Author: "Test Author",
		Pages:  100,
	}
	err := repo.Create(book)
	assert.NoError(t, err)

	// Delete it
	err = repo.Delete(book.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(book.ID)
	assert.Error(t, err)
}
