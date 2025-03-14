# Personal Notes with Go

This project is a simple REST API built with Go that allows users to manage their personal notes. It provides endpoints for creating, reading, updating, and deleting notes.

## Project Details

*   **Language:** Go
*   **Database:** In-memory (for simplicity, can be replaced with a persistent database like PostgreSQL or MySQL)
*   **Framework:** Standard Go library (no external web frameworks used)
*   **Purpose:** To demonstrate a basic REST API implementation in Go.

## Endpoints

### Notes

*   **`POST /notes`**
    *   **Description:** Creates a new note.
    *   **Request Body:**
        ```json
        {
          "title": "My Note Title",
          "content": "This is the content of my note."
        }
        ```
    *   **Response:**
        *   `201 Created`: Note created successfully.
        *   `400 Bad Request`: Invalid request body.
        *   `500 Internal Server Error`: Server error.
*   **`GET /notes`**
    *   **Description:** Retrieves all notes.
    *   **Response:**
        *   `200 OK`: Returns a list of notes.
        *   `500 Internal Server Error`: Server error.
*   **`GET /notes/{id}`**
    *   **Description:** Retrieves a specific note by ID.
    *   **Response:**
        *   `200 OK`: Returns the requested note.
        *   `404 Not Found`: Note not found.
        *   `500 Internal Server Error`: Server error.
*   **`PUT /notes/{id}`**
    *   **Description:** Updates a specific note by ID.
    *   **Request Body:**
        ```json
        {
          "title": "Updated Note Title",
          "content": "This is the updated content."
        }
        
    *   **Response:**
        *   `200 OK`: Note updated successfully.
        *   `400 Bad Request`: Invalid request body.
        *   `404 Not Found`: Note not found.
        *   `500 Internal Server Error`: Server error.
*   **`DELETE /notes/{id}`**
    *   **Description:** Deletes a specific note by ID.
    *   **Response:**
        *   `204 No Content`: Note deleted successfully.
        *   `404 Not Found`: Note not found.
        *   `500 Internal Server Error`: Server error.
    
### Categories

*   **`POST /categories`**
    *   **Description:** Creates a new category.
    *   **Request Body:**
        ```json
        {
          "name": "Work"
        }
        ```
    *   **Response:**
        *   `201 Created`: Category created successfully.
        *   `400 Bad Request`: Invalid request body.
        *   `500 Internal Server Error`: Server error.
*   **`GET /categories`**
    *   **Description:** Retrieves all categories.
    *   **Response:**
        *   `200 OK`: Returns a list of categories.
        *   `500 Internal Server Error`: Server error.
*   **`GET /categories/{id}`**
    *   **Description:** Retrieves a specific category by ID.
    *   **Response:**
        *   `200 OK`: Returns the requested category.
        *   `404 Not Found`: Category not found.
        *   `500 Internal Server Error`: Server error.
*   **`PUT /categories/{id}`**
    *   **Description:** Updates a specific category by ID.
    *   **Request Body:**
        ```json
        {
          "name": "Updated Work"
        }
        ```
    *   **Response:**
        *   `200 OK`: Category updated successfully.
        *   `400 Bad Request`: Invalid request body.
        *   `404 Not Found`: Category not found.
        *   `500 Internal Server Error`: Server error.
*   **`DELETE /categories/{id}`**
    *   **Description:** Deletes a specific category by ID.
    *   **Response:**
        *   `204 No Content`: Category deleted successfully.
        *   `404 Not Found`: Category not found.
        *   `500 Internal Server Error`: Server error.

## How to Run

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd personal-notes-with-go
    ```
2.  **Run the application:**
    ```bash
    go run main.go
    ```
3.  **The API will be available at `http://localhost:8080`**

## Curl for Testing

### Notes

*   **Create a note:**
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"title": "My First Note", "content": "This is the content of my first note."}' http://localhost:8080/notes
    ```
*   **Get all notes:**
    ```bash
    curl http://localhost:8080/notes
    ```
*   **Get a note by ID (replace {id} with the actual ID):**
    ```bash
    curl http://localhost:8080/notes/{id}
    ```
*   **Update a note (replace {id} with the actual ID):**
    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"title": "Updated Note", "content": "This is the updated content."}' http://localhost:8080/notes/{id}
    ```
*   **Delete a note (replace {id} with the actual ID):**
    ```bash
    curl -X DELETE http://localhost:8080/notes/{id}
    ```

### Categories

*   **Create a category:**
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"name": "Personal"}' http://localhost:8080/categories
    ```
*   **Get all categories:**
    ```bash
    curl http://localhost:8080/categories
    ```
*   **Get a category by ID (replace {id} with the actual ID):**
    ```bash
    curl http://localhost:8080/categories/{id}
    ```
*   **Update a category (replace {id} with the actual ID):**
    curl -X PUT -H "Content-Type: application/json" -d '{"name": "Updated Personal"}' http://localhost:8080/categories/{id}
    ```
*   **Delete a category (replace {id} with the actual ID):**
    ```bash
    curl -X DELETE http://localhost:8080/categories/{id}
    ```

## License

This project is licensed under a **No Commercial Use License**.

This means that while you are free to use, modify, and distribute this project for personal or educational purposes, you are strictly prohibited from selling it or using it for any commercial gain.

