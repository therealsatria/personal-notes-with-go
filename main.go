package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    _ "github.com/mattn/go-sqlite3"
    "net/http"
)

type Category struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type Note struct {
    ID         string `json:"id"`
    Subject    string `json:"subject"`
    Content    string `json:"content"`
    Priority   string `json:"priority"`
    Tags       string `json:"tags"`
    CategoryID string `json:"category_id"`
}

func main() {
    // Inisialisasi database SQLite
    db, err := sql.Open("sqlite3", "./db.sqlite3")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Buat tabel categories
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS categories (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL
    )`)
    if err != nil {
        panic(err)
    }

    // Buat tabel notes
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS notes (
        id TEXT PRIMARY KEY,
        subject TEXT NOT NULL,
        content TEXT,
        priority TEXT,
        tags TEXT,
        category_id TEXT,
        FOREIGN KEY (category_id) REFERENCES categories(id)
    )`)
    if err != nil {
        panic(err)
    }

    // Inisialisasi Gin
    r := gin.Default()

    // Endpoint: Tambah kategori
    r.POST("/categories", func(c *gin.Context) {
        var cat Category
        if err := c.BindJSON(&cat); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        cat.ID = uuid.New().String()
        _, err := db.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", cat.ID, cat.Name)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, cat)
    })

    // Endpoint: Ambil semua kategori
    r.GET("/categories", func(c *gin.Context) {
        rows, err := db.Query("SELECT id, name FROM categories")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()

        var categories []Category
        for rows.Next() {
            var cat Category
            if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            categories = append(categories, cat)
        }
        c.JSON(http.StatusOK, categories)
    })

    // Endpoint: Update kategori (PUT)
    r.PUT("/categories/:id", func(c *gin.Context) {
        id := c.Param("id")
        var cat Category
        if err := c.BindJSON(&cat); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        result, err := db.Exec("UPDATE categories SET name = ? WHERE id = ?", cat.Name, id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        rowsAffected, _ := result.RowsAffected()
        if rowsAffected == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
            return
        }

        cat.ID = id
        c.JSON(http.StatusOK, cat)
    })

    // Endpoint: Hapus kategori (DELETE)
    r.DELETE("/categories/:id", func(c *gin.Context) {
        id := c.Param("id")
        result, err := db.Exec("DELETE FROM categories WHERE id = ?", id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        rowsAffected, _ := result.RowsAffected()
        if rowsAffected == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
    })

    // Endpoint: Tambah note
    r.POST("/notes", func(c *gin.Context) {
        var note Note
        if err := c.BindJSON(¬e); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        note.ID = uuid.New().String()
        _, err := db.Exec("INSERT INTO notes (id, subject, content, priority, tags, category_id) VALUES (?, ?, ?, ?, ?, ?)",
            note.ID, note.Subject, note.Content, note.Priority, note.Tags, note.CategoryID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, note)
    })

    // Endpoint: Ambil semua notes (dengan filter priority dan category)
    r.GET("/notes", func(c *gin.Context) {
        priority := c.Query("priority")     // Filter berdasarkan priority
        categoryID := c.Query("category_id") // Filter berdasarkan category_id

        query := "SELECT id, subject, content, priority, tags, category_id FROM notes WHERE 1=1"
        var args []interface{}

        if priority != "" {
            query += " AND priority = ?"
            args = append(args, priority)
        }
        if categoryID != "" {
            query += " AND category_id = ?"
            args = append(args, categoryID)
        }

        rows, err := db.Query(query, args...)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer rows.Close()

        var notes []Note
        for rows.Next() {
            var n Note
            if err := rows.Scan(&n.ID, &n.Subject, &n.Content, &n.Priority, &n.Tags, &n.CategoryID); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            notes = append(notes, n)
        }
        c.JSON(http.StatusOK, notes)
    })

    // Endpoint: Update note (PUT)
    r.PUT("/notes/:id", func(c *gin.Context) {
        id := c.Param("id")
        var note Note
        if err := c.BindJSON(¬e); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        result, err := db.Exec("UPDATE notes SET subject = ?, content = ?, priority = ?, tags = ?, category_id = ? WHERE id = ?",
            note.Subject, note.Content, note.Priority, note.Tags, note.CategoryID, id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        rowsAffected, _ := result.RowsAffected()
        if rowsAffected == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
            return
        }

        note.ID = id
        c.JSON(http.StatusOK, note)
    })

    // Endpoint: Hapus note (DELETE)
    r.DELETE("/notes/:id", func(c *gin.Context) {
        id := c.Param("id")
        result, err := db.Exec("DELETE FROM notes WHERE id = ?", id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        rowsAffected, _ := result.RowsAffected()
        if rowsAffected == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
    })

    // Jalankan server
    r.Run(":8080")
}