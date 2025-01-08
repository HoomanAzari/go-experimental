package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})

	// Start the server on port 8080
	fmt.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
