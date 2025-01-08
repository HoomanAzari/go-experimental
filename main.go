package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Respond with a simple health check message
		w.WriteHeader(http.StatusOK) // HTTP 200 OK
		fmt.Fprintf(w, "Server is healthy!")
	})

	// Starting the server on port 8080
	fmt.Println("Health check server starting on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
