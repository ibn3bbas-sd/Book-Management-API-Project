package main

// Sample code for a simple book management API in Go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	PublicationYear	int    `json:"published_year"`
}

var (
	books      = make(map[int]Book)
	mutex      = &sync.Mutex{}
)

func main() {