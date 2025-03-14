package main

import (
	"log"
	"personal-notes-with-go/database"
	"personal-notes-with-go/handlers"
	"personal-notes-with-go/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi database
	db, err := database.InitDB("./db.sqlite3")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Inisialisasi repository
	categoryRepo := repositories.NewCategoryRepository(db)
	noteRepo := repositories.NewNoteRepository(db)

	// Inisialisasi Gin
	r := gin.Default()

	// Inisialisasi handler
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	noteHandler := handlers.NewNoteHandler(noteRepo)

	// Routing
	categoryGroup := r.Group("/categories")
	{
		categoryGroup.POST("", categoryHandler.CreateCategory)
		categoryGroup.GET("", categoryHandler.GetCategories)
		categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
		categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	noteGroup := r.Group("/notes")
	{
		noteGroup.POST("", noteHandler.CreateNote)
		noteGroup.GET("", noteHandler.GetNotes)
		noteGroup.PUT("/:id", noteHandler.UpdateNote)
		noteGroup.DELETE("/:id", noteHandler.DeleteNote)
	}

	// Jalankan server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
