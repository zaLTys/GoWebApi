// @title Books API
// @version 1.0
// @description Example API with Gin, GORM, and Swagger docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email you@example.com
// @license.name MIT
// @host localhost:8080
// @BasePath /

package main

import (
	"log"
	"net/http"

	_ "books-api/docs" // 👈 replace with your module name

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Auto-migrate creates table 'books'
	if err := db.AutoMigrate(&Book{}); err != nil {
		log.Fatal("failed to migrate database")
	}
}

// createBook godoc
// @Summary      Create a new book
// @Description  Adds a new book to the database
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book body Book true "Book data"
// @Success      201 {object} Book
// @Failure      400 {object} map[string]string
// @Router       /books [post]
func createBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if book.Color != nil && !book.Color.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid color"})
		return
	}
	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// listBooks godoc
// @Summary      List all books
// @Description  Get all books
// @Tags         books
// @Produce      json
// @Success      200 {array} Book
// @Router       /books [get]
func listBooks(c *gin.Context) {
	var books []Book
	db.Find(&books)
	c.JSON(http.StatusOK, books)
}

// getBook godoc
// @Summary      Get a book by ID
// @Description  Returns a single book
// @Tags         books
// @Produce      json
// @Param        id path int true "Book ID"
// @Success      200 {object} Book
// @Failure      404 {object} map[string]string
// @Router       /books/{id} [get]
func getBook(c *gin.Context) {
	var book Book
	if err := db.First(&book, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// updateBook godoc
// @Summary      Update a book
// @Description  Updates book fields by ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id path int true "Book ID"
// @Param        book body Book true "Book data"
// @Success      200 {object} Book
// @Failure      400 {object} map[string]string
// @Router       /books/{id} [put]
func updateBook(c *gin.Context) {
	var book Book
	if err := db.First(&book, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	var input Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Color != nil && !input.Color.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid color"})
		return
	}

	// Use a helper method to update fields
	book.Update(input)

	if err := db.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// deleteBook godoc
// @Summary      Delete a book
// @Description  Deletes a book by ID
// @Tags         books
// @Produce      json
// @Param        id path int true "Book ID"
// @Success      200 {object} map[string]string
// @Router       /books/{id} [delete]
func deleteBook(c *gin.Context) {
	if err := db.Delete(&Book{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func main() {
	initDB()

	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/books", createBook)
	r.GET("/books", listBooks)
	r.GET("/books/:id", getBook)
	r.PUT("/books/:id", updateBook)
	r.DELETE("/books/:id", deleteBook)

	r.Run(":8080")
}
