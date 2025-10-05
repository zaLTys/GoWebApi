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
	"os"

	_ "books-api/docs"
	"books-api/app/controller"
	"books-api/app/migrations"
	"books-api/app/repository"
	"books-api/app/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Run migrations
	migrationManager := migrations.NewMigrationManager()
	if err := migrationManager.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize layers
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookController := controller.NewBookController(bookService)

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	setupRoutes(r, bookController)

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDB initializes the database connection
func initDB() (*gorm.DB, error) {
	log.Println("Initializing database connection...")
	
	db, err := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// setupRoutes configures all the API routes
func setupRoutes(r *gin.Engine, bookController *controller.BookController) {
	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Book routes
	bookRoutes := r.Group("/books")
	{
		bookRoutes.POST("", bookController.CreateBook)
		bookRoutes.GET("", bookController.ListBooks)
		bookRoutes.GET("/:id", bookController.GetBook)
		bookRoutes.PUT("/:id", bookController.UpdateBook)
		bookRoutes.DELETE("/:id", bookController.DeleteBook)
	}

	log.Println("Routes configured successfully")
}