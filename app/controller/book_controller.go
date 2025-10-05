package controller

import (
	"books-api/app/service"
	"books-api/app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookController handles HTTP requests for book operations
type BookController struct {
	bookService service.BookService
}

// NewBookController creates a new instance of book controller
func NewBookController(bookService service.BookService) *BookController {
	return &BookController{
		bookService: bookService,
	}
}

// CreateBook godoc
// @Summary      Create a new book
// @Description  Adds a new book to the database
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book body models.Book true "Book data"
// @Success      201 {object} models.Book
// @Failure      400 {object} map[string]string
// @Router       /books [post]
func (ctrl *BookController) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.bookService.CreateBook(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// ListBooks godoc
// @Summary      List all books
// @Description  Get all books
// @Tags         books
// @Produce      json
// @Success      200 {array} models.Book
// @Router       /books [get]
func (ctrl *BookController) ListBooks(c *gin.Context) {
	books, err := ctrl.bookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetBook godoc
// @Summary      Get a book by ID
// @Description  Returns a single book
// @Tags         books
// @Produce      json
// @Param        id path int true "Book ID"
// @Success      200 {object} models.Book
// @Failure      404 {object} map[string]string
// @Router       /books/{id} [get]
func (ctrl *BookController) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	book, err := ctrl.bookService.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// UpdateBook godoc
// @Summary      Update a book
// @Description  Updates book fields by ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id path int true "Book ID"
// @Param        book body models.Book true "Book data"
// @Success      200 {object} models.Book
// @Failure      400 {object} map[string]string
// @Router       /books/{id} [put]
func (ctrl *BookController) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	var updateData models.Book
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := ctrl.bookService.UpdateBook(uint(id), updateData)
	if err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary      Delete a book
// @Description  Deletes a book by ID
// @Tags         books
// @Produce      json
// @Param        id path int true "Book ID"
// @Success      200 {object} map[string]string
// @Router       /books/{id} [delete]
func (ctrl *BookController) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	if err := ctrl.bookService.DeleteBook(uint(id)); err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}
