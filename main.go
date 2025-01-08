package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Root Handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the homepage!")
	})

	// About Hnadler function
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is the About page!")
	})

	// Contact Handler function
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Contact me at veryhouman@gmail.com")
	})

	// Start the server on port 8080
	fmt.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
