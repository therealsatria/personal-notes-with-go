package main

import (
	"log"
	"net/http"
	"personal-notes-with-go/database"
	"personal-notes-with-go/handlers"
	"personal-notes-with-go/repositories"
	"personal-notes-with-go/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Middleware to check if encryption is valid before allowing data modification
func requireValidEncryption() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsEncryptionValid() {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Encryption system is not properly initialized. Data modification is disabled for security reasons.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	// Initialize encryption
	if err := utils.InitEncryption(); err != nil {
		log.Printf("WARNING: Failed to initialize encryption: %v", err)
		log.Printf("Data modification will be disabled for security reasons.")
		// We continue execution but with encryption marked as invalid
	}

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

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Serve static files for frontend
	r.Static("/frontend", "./frontend")

	// Serve the SPA
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/frontend")
	})

	// Inisialisasi handler
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	noteHandler := handlers.NewNoteHandler(noteRepo)
	keyHandler := handlers.NewKeyHandler()
	encryptionHandler := handlers.NewEncryptionHandler()

	// Encryption status endpoint
	r.GET("/encryption/status", encryptionHandler.GetStatus)

	// Routing with encryption validation middleware for data modification endpoints
	categoryGroup := r.Group("/categories")
	{
		categoryGroup.POST("", requireValidEncryption(), categoryHandler.CreateCategory)
		categoryGroup.GET("", categoryHandler.GetCategories)
		categoryGroup.PUT("/:id", requireValidEncryption(), categoryHandler.UpdateCategory)
		categoryGroup.DELETE("/:id", requireValidEncryption(), categoryHandler.DeleteCategory)
	}

	noteGroup := r.Group("/notes")
	{
		noteGroup.POST("", requireValidEncryption(), noteHandler.CreateNote)
		noteGroup.GET("", noteHandler.GetNotes)
		noteGroup.PUT("/:id", requireValidEncryption(), noteHandler.UpdateNote)
		noteGroup.DELETE("/:id", requireValidEncryption(), noteHandler.DeleteNote)
	}

	// Key generation endpoint
	r.POST("/generate-key", keyHandler.GenerateKey)

	// Jalankan server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
