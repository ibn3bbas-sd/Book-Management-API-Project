package main

// Sample code for a simple book management API in Go

import (
	"encoding/json" // encoding/json package for JSON encoding and decoding
	"fmt"           // package for formatted I/O
	"log"           // log package for logging
	"net/http"      // net/http package for handling HTTP requests
	"strconv"       // strconv package for converting strings to integers
	"strings"       // strings package for string manipulation
	"sync"          // sync package for synchronization primitives
)

// JSON syntax for the Book struct `json:"id"`

type Book struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"published_year"`
}

var (
	books      = make(map[int]Book)
	mutex      = &sync.Mutex{}
	nextBookID = 1 // counter to generate new unique IDs
)

func main() {
	http.HandleFunc("/books", getBooks) // Handle GET requests to /books
	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBook(w, r) // Handle GET requests to /books/{id}
		case http.MethodPost:
			createBook(w, r) // Handle POST requests to /books
		case http.MethodPut:
			updateBook(w, r) // Handle PUT requests to /books/{id}
		case http.MethodDelete:
			deleteBook(w, r) // Handle DELETE requests to /books/{id}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // Handle unsupported methods
		}
	})

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Function for our CRUD API
func getBooks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock() // Ensure that the mutex is unlocked after the function completes

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get Book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock() // Ensure that the mutex is unlocked after the function completes

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, exists := books[id]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock() // Ensure that the mutex is unlocked after the function completes

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// User didn't provide an ID, so we generate one
	if book.ID == 0 {
		book.ID = nextBookID
		nextBookID++
	} else {
		// If ID is provided, make sure it doesn't already exist
		if _, exists := books[book.ID]; exists {
			http.Error(w, "Book ID already exists", http.StatusBadRequest)
			return
		}
	}

	books[book.ID] = book
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// Update an existing book
func updateBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock() // Ensure that the mutex is unlocked after the function completes

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updatedBook.ID != id {
		http.Error(w, "Book ID in URL and body do not match", http.StatusBadRequest)
		return
	}

	if _, exists := books[id]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	books[id] = updatedBook
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

// Delete an existing book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock() // Ensure that the mutex is unlocked after the function completes

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if _, exists := books[id]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	delete(books, id)
	w.WriteHeader(http.StatusNoContent) // No content to return
}
