package main

// Sample code for a simple book management API in Go

import (
	"encoding/json" // encoding/json package for JSON encoding and decoding
	// fmt package for formatted I/O
	"net/http" // net/http package for handling HTTP requests
	"strconv"  // strconv package for converting strings to integers
	"strings"  // strings package for string manipulation
	"sync"     // sync package for synchronization primitives
)

// JSON syntax for the Book struct `json:"id"`

type Book struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"published_year"`
}

var (
	books = make(map[int]Book)
	mutex = &sync.Mutex{}
)

func main() {

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
	if err != nil /*|| id <= 0*/ {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	book, exists := books[id]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

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
	// User ID did't provide an ID, so we generate one
	if book.ID <= 0 || book.Title == "" || book.Author == "" || book.PublicationYear <= 0 { // Need edit
		http.Error(w, "Invalid book data", http.StatusBadRequest)
		return
	}

	books[book.ID] = book
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Update an existing book
func updateBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock() // Ensure that the mutex is unlocked after the function completes

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil /*|| id <= 0*/ {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var updateBook Book
	if err := json.NewDecoder(r.Body).Decode(&updateBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateBook.ID != id || updateBook.Title == "" || updateBook.Author == "" || updateBook.PublicationYear <= 0 { // Need edit
		http.Error(w, "Invalid book data", http.StatusBadRequest)
		return
	}

	books[id] = updateBook
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateBook)
}

// Delete an existing book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock() // Ensure that the mutex is unlocked after the function completes
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil /*|| id <= 0*/ {
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
